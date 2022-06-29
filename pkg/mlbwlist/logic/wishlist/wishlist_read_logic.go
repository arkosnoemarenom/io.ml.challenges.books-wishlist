package wishlist

import "io.ml.challenges/io.ml.challenges.books-wishlist/pkg/mlbwlist/models"

func (l *logic) Read(userid string, wishlistid string) (*models.Wishlist, error) {

	item, err := l.db.Select(userid, wishlistid)
	if err != nil {
		return nil, err
	}

	return item, nil
}