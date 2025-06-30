package hardcover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"web_pull/sources"
)

type HardcoverScraper struct{}

const hardcoverAPIURL = "https://api.hardcover.app/graphql"

// Struct for the GraphQL request body
type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

// Struct to parse the book ID from the response
type BookIDResponse struct {
	Data struct {
		Book struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"book"`
	} `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// Struct to parse the reviews from the response
type ReviewsResponse struct {
	Data struct {
		Reviews struct {
			Nodes []struct {
				ID      string `json:"id"`
				Rating  *int   `json:"rating"` // Use pointer to handle null ratings
				Content string `json:"content"`
				User    struct {
					ID       string `json:"id"`
					Username string `json:"username"`
				} `json:"user"`
			} `json:"nodes"`
		} `json:"reviews"`
	} `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// Generic GraphQL error struct
type GraphQLError struct {
	Message string `json:"message"`
}

// Helper function to make GraphQL requests
func makeGraphQLRequest(reqBody []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", hardcoverAPIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func (h HardcoverScraper) GetBookIDByISBN(isbn string) (string, error) {
	query := `
    query GetBookByIsbn($isbn: String!) {
      book(isbn: $isbn) {
        id
        title
      }
    }`
	variables := map[string]interface{}{"isbn": isbn}
	reqBody, _ := json.Marshal(GraphQLRequest{Query: query, Variables: variables})

	respBody, err := makeGraphQLRequest(reqBody)
	if err != nil {
		return "", err
	}

	var parsedResp BookIDResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal book ID response: %w", err)
	}

	if len(parsedResp.Errors) > 0 {
		return "", fmt.Errorf("API error fetching book ID: %s", parsedResp.Errors[0].Message)
	}
	if parsedResp.Data.Book.ID == "" {
		return "", fmt.Errorf("no book found with ISBN: %s", isbn)
	}

	return parsedResp.Data.Book.ID, nil
}

func (h HardcoverScraper) GetReviews(bookID string) ([]sources.ItemizedReview, error) {
	query := `
    query GetReviewsForBook($bookId: ID!) {
      reviews(where: { book_id: { _eq: $bookId } }) {
        nodes {
          id
          rating
          content
          user {
            id
            username
          }
        }
      }
    }`
	variables := map[string]interface{}{"bookId": bookID}
	reqBody, _ := json.Marshal(GraphQLRequest{Query: query, Variables: variables})

	respBody, err := makeGraphQLRequest(reqBody)
	if err != nil {
		return nil, err
	}

	var parsedResp ReviewsResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reviews response: %w", err)
	}

	if len(parsedResp.Errors) > 0 {
		return nil, fmt.Errorf("API error fetching reviews: %s", parsedResp.Errors[0].Message)
	}

	var itemizedReviews []sources.ItemizedReview
	for _, reviewNode := range parsedResp.Data.Reviews.Nodes {
		itemizedReviews = append(itemizedReviews, sources.ItemizedReview{
			UniqueIdentifier: reviewNode.ID,
			Rating:           reviewNode.Rating,
			Text:             reviewNode.Content,
			UserID:           reviewNode.User.ID,
			Username:         reviewNode.User.Username,
		})
	}

	return itemizedReviews, nil
}

/*
func main() {
	// Example using "Dune" by Frank Herbert
	testISBN := "9780441013593"

	log.Printf("Attempting to find book with ISBN: %s", testISBN)
	bookID, err := getBookIDByISBN(testISBN)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Found book with Hardcover ID: %s. Fetching reviews...", bookID)
	reviews, err := getReviewsForBook(bookID)
	if err != nil {
		log.Fatalf("Error fetching reviews: %v", err)
	}

	if len(reviews) > 0 {
		log.Printf("Successfully fetched %d review(s).", len(reviews))
		// Print the full list of reviews as an example
		prettyJSON, err := json.MarshalIndent(reviews, "", "    ")
		if err != nil {
			log.Fatalf("Failed to generate JSON: %v", err)
		}
		fmt.Println(string(prettyJSON))
	} else {
		log.Println("No reviews found for this book.")
	}
}
*/
