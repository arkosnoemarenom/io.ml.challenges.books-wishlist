package wishlist

import (
	"fmt"

	"io.ml.challenges/io.ml.challenges.books-wishlist/pkg/mlbwlist/models"
)

func (c *connector) Insert(wishlistModel models.Wishlist) error {

	sql := fmt.Sprintf("INSERT INTO books_wishlists(_id, spec_user_id, spec_name, spec_description) VALUES('%s', '%s', '%s', '%s')", wishlistModel.ID, wishlistModel.Spec.User, wishlistModel.Spec.Name, wishlistModel.Spec.Description)
	c.logging.Info(sql)
	_, err := c.db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
