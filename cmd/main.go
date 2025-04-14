package main

import (
	"avito2/internal/db"
	"avito2/internal/handler_manager"
	"avito2/internal/middleware"
	"avito2/internal/repository"
	"avito2/internal/service"
	"avito2/internal/utils"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	httpPort := os.Getenv("SERVER_PORT")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer database.GetPool(ctx).Close()

	repo := repository.NewRepository(database)
	svc := service.NewService(repo)
	jwtGen := &utils.JWTGen{}
	hm := handler_manager.NewHandlerManager(svc, jwtGen)

	r := mux.NewRouter()
	r.Handle("/pvz", middleware.AuthMiddleware(http.HandlerFunc(hm.Pvz)))
	r.Handle("/pvz/{pvzId}/close_last_reception", middleware.AuthMiddleware(http.HandlerFunc(hm.CloseLastReception)))
	r.Handle("/pvz/{pvzId}/delete_last_product", middleware.AuthMiddleware(http.HandlerFunc(hm.DeleteLastProduct)))
	r.Handle("/receptions", middleware.AuthMiddleware(http.HandlerFunc(hm.CreateReception)))
	r.Handle("/products", middleware.AuthMiddleware(http.HandlerFunc(hm.AddProduct)))
	r.HandleFunc("/dummyLogin", hm.DummyLogin)

	log.Println("http serer start listening on port:", httpPort)
	err = http.ListenAndServe(":"+httpPort, r)

	if err != nil {
		log.Fatal(err)
	}
}
