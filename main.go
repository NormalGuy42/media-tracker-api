package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Movie struct {
	ID          int       `json:"id,omitempty"`
	MovieID     int       `json:"movieID"`
	Title       string    `json:"title,omitempty"`
	Poster      string    `json:"poster,omitempty"`
	Director    *[]string `json:"director,omitempty"`
	Cast        []string  `json:"cast,omitempty"`
	Review      string    `json:"review,omitempty"`
	Genres      []string  `json:"genres,omitempty"`
	Watched     bool      `json:"watched,omitempty"`
	PlanToWatch bool      `json:"planToWatch,omitempty"`
	Released    *string   `json:"released,omitempty"`
	Created_At  time.Time `json:"created_at"`
}

var client *supabase.Client

func main() {
	app := fiber.New()
	fmt.Println("Hello World!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api_key := os.Getenv("SUPABASE_API_KEY")
	project_ref := os.Getenv("SUPABASE_PROJECT_REF")

	client, err := supabase.NewClient(api_key, project_ref, nil)

	if err != nil {
		fmt.Println("failed in init supa client: ", err)
		return
	}

	PORT := os.Getenv("PORT")

	//Movies
	app.Get("/api/movies", func(c *fiber.Ctx) error { return getMovies(c, client) })
	app.Post("/api/movies", func(c *fiber.Ctx) error { return createMovie(c, client) })
	app.Get("/api/movies/:id", func(c *fiber.Ctx) error { return getMovieByID(c, client) })
	app.Patch("/api/movies/:id", func(c *fiber.Ctx) error { return updateMovie(c, client) })
	app.Delete("/api/movies/:id", func(c *fiber.Ctx) error { return deleteMovie(c, client) })

	//Books
	// app.Get("/api/books",getBooks)
	// app.Get("/api/books/:id",getBookByID)
	// app.Post("/api/books/:id",createBook)
	// app.Patch("/api/books/:id",updateBook)
	// app.Delete("/api/books/:id",deleteBook)

	// //TV Shows
	// app.Get("/api/shows",getShows)
	// app.Get("/api/shows/:id",getShowByID)
	// app.Post("/api/shows/:id",createShow)
	// app.Patch("/api/shows/:id",updateShow)
	// app.Delete("/api/shows/:id",deleteShow)

	// //Game
	// app.Get("/api/games",getGame)
	// app.Get("/api/games/:id",getGameByID)
	// app.Post("/api/games/:id",createGame)
	// app.Patch("/api/games/:id",updateGame)
	// app.Delete("/api/games/:id",deleteGame)

	// if PORT == "" {
	// 	PORT = "5000"
	// }
	log.Fatal(app.Listen(":" + PORT))
}

func getMovies(c *fiber.Ctx, client *supabase.Client) error {

	var movies []Movie
	_, err := client.From("movies").Select("*", "", false).ExecuteTo(&movies)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(movies)
}

func getMovieByID(c *fiber.Ctx, client *supabase.Client) error {
	id := c.Params("id")
	var movies []Movie
	_, err := client.From("movies").Select("*", "", false).Eq("id", id).ExecuteTo(&movies)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if len(movies) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Movie not found"})
	}
	return c.JSON(movies[0])
}

func createMovie(c *fiber.Ctx, client *supabase.Client) error {
	var movie Movie
	if err := c.BodyParser(&movie); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	data, _, err := client.From("movies").Insert(movie, true, "", "", "").Execute()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Safely parse inserted row(s)
	var inserted []Movie
	if err := json.Unmarshal(data, &inserted); err != nil {
		var single Movie
		if err := json.Unmarshal(data, &single); err == nil {
			inserted = append(inserted, single)
		}
	}

	return c.JSON(inserted)
}

func updateMovie(c *fiber.Ctx, client *supabase.Client) error {
	id := c.Params("id")
	var movie Movie
	if err := c.BodyParser(&movie); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	data, _, err := client.From("movies").Update(movie, "", "").Eq("id", id).Execute()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var updated []Movie
	_ = json.Unmarshal(data, &updated)
	return c.JSON(updated)
}

func deleteMovie(c *fiber.Ctx, client *supabase.Client) error {
	id := c.Params("id")
	_, _, err := client.From("movies").Delete("", "").Eq("id", id).Execute()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(204).JSON(fiber.Map{"message": "Movie Deleted"})
}

//Books

// func createBook(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func getBooks(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func getBookByID(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func updateBook(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func deleteBook(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }

// //TV shows

// func createShow(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func getShows(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func getShowByID(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func updateShow(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func deleteShow(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// //Games

// func createGame(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func getGame(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func getGameByID(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func updateGame(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
// func deleteGame(c *fiber.Ctx) error {
// 	id := c.Params("id")
// }
