package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	supabase "github.com/lengzuo/supa"
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

	conf := supabase.Config{
		// Your project api key, you can use either `anon` or `service_role`.
		// but i will suggest you to use `service_role` as your api key and keep it secret.
		ApiKey: os.Getenv("SUPABASE_API_KEY"),
		// Retrieve your project ref from project url
		// eg: https://this-your-project-ref.supabase.co
		ProjectRef: os.Getenv("SUPABASE_PROJECT_REF"),
		// Set it `false` in production to avoid extra log print.
		Debug: true,
	}

	client, err := supabase.New(conf)

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var movies []Movie
	if err := client.DB.From("movies").Select("*").Execute(ctx, &movies); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(movies)
}

func getMovieByID(c *fiber.Ctx, client *supabase.Client) error {
	id := c.Params("id")
	var movies []Movie

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.DB.From("movies").Select("*").Eq("id", id).Execute(ctx, &movies)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var inserted []Movie

	if err := client.DB.From("movies").Insert(movie).Execute(ctx, &inserted); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(inserted)
}

func updateMovie(c *fiber.Ctx, client *supabase.Client) error {
	id := c.Params("id")
	movie := new(Movie)

	if err := c.BodyParser(movie); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var updatedMovie []Movie

	err := client.DB.From("movies").Update(movie).Eq("id", id).Execute(ctx, &updatedMovie)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updatedMovie)
}
func deleteMovie(c *fiber.Ctx, client *supabase.Client) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.DB.From("movies").Delete().Eq("id", id).Execute(ctx, nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Movie Deleted"})
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
