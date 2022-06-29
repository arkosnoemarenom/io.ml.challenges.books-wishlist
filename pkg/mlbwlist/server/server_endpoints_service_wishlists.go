package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"io.ml.challenges/io.ml.challenges.books-wishlist/pkg/mlbwlist/models"
)

func (s *serve) configureWishlistsServiceEndpoints() {

	createWishlistEndpointPath := fmt.Sprintf("%s%s", apiversion, wishListsEndpointPath)
	s.router.HandleFunc(createWishlistEndpointPath, s.manageCreateWishlistsRequest).
		Methods(http.MethodPost)

	readWishlistEndpointPath := fmt.Sprintf("%s%s", apiversion, wishListEndpointPath)
	s.router.HandleFunc(readWishlistEndpointPath, s.manageReadWishlistRequest).
		Methods(http.MethodGet)

	listWishlistEndpointPath := fmt.Sprintf("%s%s", apiversion, wishListsEndpointPath)
	s.router.HandleFunc(listWishlistEndpointPath, s.manageListWishlistsRequest).
		Methods(http.MethodGet)
}

func (s *serve) manageCreateWishlistsRequest(w http.ResponseWriter, r *http.Request) {

	userid := mux.Vars(r)["user"]
	wishlistModel := models.Wishlist{}
	json.NewDecoder(r.Body).Decode(&wishlistModel)
	if wishlistModel.Spec == nil {
		wishlistModel.Spec = &models.WishlistSpec{}
	}
	wishlistModel.Spec.User = userid
	model, err := s.wishlistLogic.Create(wishlistModel)
	if err != nil {
		s.httpServiceErrorManagement(w, err.Error())
		return
	}

	s.httpJsonResponseManagement(w, model, http.StatusOK)
}

func (s *serve) manageReadWishlistRequest(w http.ResponseWriter, r *http.Request) {

	muxvars := mux.Vars(r)
	userid := muxvars["user"]
	wishlistid := muxvars["wishlist"]
	item, err := s.wishlistLogic.Read(userid, wishlistid)
	if err != nil {
		s.httpServiceErrorManagement(w, err.Error())
		return
	}

	s.httpJsonResponseManagement(w, item, http.StatusOK)
}

func (s *serve) manageListWishlistsRequest(w http.ResponseWriter, r *http.Request) {

	userid := mux.Vars(r)["user"]
	list, err := s.wishlistLogic.List(userid)
	if err != nil {
		s.httpServiceErrorManagement(w, err.Error())
		return
	}

	s.httpJsonResponseManagement(w, list, http.StatusOK)
}
