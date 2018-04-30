package gomoney

import (
	"fmt"

	goerror "github.com/joaosoft/go-error/app"
	"github.com/satori/go.uuid"
)

// iStorageDB ...
type iStorageDB interface {
	getSession(userID uuid.UUID, token string) (*session, *goerror.ErrorData)
	getSessions(userID uuid.UUID) ([]*session, *goerror.ErrorData)
	createSession(newSession *session) (*session, *goerror.ErrorData)
	deleteSession(userID uuid.UUID, token string) *goerror.ErrorData
	deleteSessions(userID uuid.UUID) *goerror.ErrorData

	getUsers() ([]*user, *goerror.ErrorData)
	getUser(userID uuid.UUID) (*user, *goerror.ErrorData)
	getUserByEmail(email string) (*user, *goerror.ErrorData)
	createUser(newUser *user) (*user, *goerror.ErrorData)
	updateUser(updUser *user) (*user, *goerror.ErrorData)
	deleteUser(userID uuid.UUID) *goerror.ErrorData

	getWallets(userID uuid.UUID) ([]*wallet, *goerror.ErrorData)
	getWallet(userID uuid.UUID, walletID uuid.UUID) (*wallet, *goerror.ErrorData)
	createWallets(newWallets []*wallet) ([]*wallet, *goerror.ErrorData)
	updateWallet(updCategory *wallet) (*wallet, *goerror.ErrorData)
	deleteWallet(userID uuid.UUID, walletID uuid.UUID) *goerror.ErrorData

	getImages(userID uuid.UUID) ([]*image, *goerror.ErrorData)
	getImage(userID uuid.UUID, imageID uuid.UUID) (*image, *goerror.ErrorData)
	createImage(newImage *image) (*image, *goerror.ErrorData)
	updateImage(updImage *image) (*image, *goerror.ErrorData)
	deleteImage(userID uuid.UUID, imageID uuid.UUID) *goerror.ErrorData

	getCategories(userID uuid.UUID) ([]*category, *goerror.ErrorData)
	getCategory(userID uuid.UUID, categoryID uuid.UUID) (*category, *goerror.ErrorData)
	createCategories(newCategory []*category) ([]*category, *goerror.ErrorData)
	updateCategory(updCategory *category) (*category, *goerror.ErrorData)
	deleteCategory(userID uuid.UUID, categoryID uuid.UUID) *goerror.ErrorData

	getTransactions(userID uuid.UUID) ([]*transaction, *goerror.ErrorData)
	getTransaction(userID uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID) (*transaction, *goerror.ErrorData)
	createTransactions(newTransaction []*transaction) ([]*transaction, *goerror.ErrorData)
	updateTransaction(updTransaction *transaction) (*transaction, *goerror.ErrorData)
	deleteTransaction(userID uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID) *goerror.ErrorData
}

// iStorageDropbox ...
type iStorageDropbox interface {
	upload(path string, data []byte) *goerror.ErrorData
	download(path string) ([]byte, *goerror.ErrorData)
	delete(path string) *goerror.ErrorData
}

// interactor ...
type interactor struct {
	storageDB      iStorageDB
	storageDropbox iStorageDropbox
	config         *MoneyConfig
}

// newInteractor ...
func newInteractor(storageDB iStorageDB, storageDropbox iStorageDropbox, config *MoneyConfig) *interactor {
	return &interactor{
		storageDB:      storageDB,
		storageDropbox: storageDropbox,
		config:         config,
	}
}

// getUsers ...
func (interactor *interactor) getUsers() ([]*user, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getUsers"})
	log.Info("getting users")
	if users, err := interactor.storageDB.getUsers(); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting users on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return users, nil
	}
}

// getUser ...
func (interactor *interactor) getUser(userID uuid.UUID) (*user, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getUser"})
	log.Infof("getting user %s", userID.String())
	if user, err := interactor.storageDB.getUser(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting user on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return user, nil
	}
}

// getUserByEmail ...
func (interactor *interactor) getUserByEmail(email string) (*user, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getUserByEmail"})
	log.Infof("getting user by email %s", email)
	if user, err := interactor.storageDB.getUserByEmail(email); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting user by email %s on storage database %s", email, err).ToErrorData(err)
		return nil, err
	} else {
		return user, nil
	}
}

// createUser ...
func (interactor *interactor) createUser(newUser *user) (*user, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "createUser"})

	newUser.UserID = uuid.NewV4()
	passwordToken, err := generateToken(authentication, []byte(newUser.Password))
	if err != nil {
		newErr := goerror.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when generating password token").ToErrorData(newErr)
		return nil, newErr
	}
	newUser.Token = passwordToken

	log.Infof("creating user %s", newUser.UserID.String())

	if user, err := interactor.storageDB.createUser(newUser); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating user on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return user, nil
	}
}

// updateUser ...
func (interactor *interactor) updateUser(updUser *user) (*user, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "updateUser"})
	log.Infof("updating user %s", updUser.UserID.String())

	passwordToken, err := generateToken(authentication, []byte(updUser.Password))
	if err != nil {
		newErr := goerror.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when generating password token").ToErrorData(newErr)
		return nil, newErr
	}
	updUser.Token = passwordToken

	if user, err := interactor.storageDB.updateUser(updUser); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error updating user on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return user, nil
	}
}

// deleteUser ...
func (interactor *interactor) deleteUser(userID uuid.UUID) *goerror.ErrorData {
	log.WithFields(map[string]interface{}{"method": "deleteUser"})
	log.Infof("deleting user %s", userID.String())
	if err := interactor.storageDB.deleteUser(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting user on storage database %s", err).ToErrorData(err)
		return err
	}
	return nil
}

// getSessions ...
func (interactor *interactor) getSessions(userID uuid.UUID) ([]*session, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getSessions"})
	log.Info("getting sessions")
	if sessions, err := interactor.storageDB.getSessions(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting sessions on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return sessions, nil
	}
}

// getSession ...
func (interactor *interactor) getSession(userID uuid.UUID, token string) (*session, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getSession"})
	log.Infof("getting session with token %s", userID.String())
	if session, err := interactor.storageDB.getSession(userID, token); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting session on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return session, nil
	}
}

// createSession ...
func (interactor *interactor) createSession(newSession *session) (*session, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "createSession"})

	random := uuid.NewV4()
	newSession.SessionID = uuid.NewV4()
	token, err := generateToken(authentication, []byte(random.String()))
	if err != nil {
		newErr := goerror.NewError(err)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting sessions on storage database %s", err).ToErrorData(newErr)
		return nil, newErr
	}
	newSession.Original = random.String()
	newSession.Token = token

	log.Infof("creating session for user %s", newSession.UserID.String())

	if user, err := interactor.storageDB.createSession(newSession); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating session on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return user, nil
	}
}

// deleteSessions ...
func (interactor *interactor) deleteSessions(userID uuid.UUID) *goerror.ErrorData {
	log.WithFields(map[string]interface{}{"method": "deleteSessions"})
	log.Infof("deleting sessions of user %s", userID.String())
	if err := interactor.storageDB.deleteSessions(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting sessions of user %s on storage database %s", userID.String(), err).ToErrorData(err)
		return err
	}
	return nil
}

// deleteSession ...
func (interactor *interactor) deleteSession(userID uuid.UUID, token string) *goerror.ErrorData {
	log.WithFields(map[string]interface{}{"method": "deleteSession"})
	log.Infof("deleting session with token %s", userID.String())
	if err := interactor.storageDB.deleteSession(userID, token); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting session of user %s  with token %s on storage database %s", userID.String(), token, err).ToErrorData(err)
		return err
	}
	return nil
}

// getWallets ...
func (interactor *interactor) getWallets(userID uuid.UUID) ([]*wallet, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getWallets"})
	log.Infof("getting wallets of user %s", userID.String())
	if wallets, err := interactor.storageDB.getWallets(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting wallets on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return wallets, nil
	}
}

// getWallet ...
func (interactor *interactor) getWallet(userID uuid.UUID, walletID uuid.UUID) (*wallet, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getWallet"})
	log.Infof("getting wallet %s of user %s", walletID.String(), userID.String())
	if wallet, err := interactor.storageDB.getWallet(userID, walletID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting wallet on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return wallet, nil
	}
}

// createWallets ...
func (interactor *interactor) createWallets(newWallets []*wallet) ([]*wallet, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "createWallets"})

	log.Info("creating wallets")
	for _, wallet := range newWallets {
		wallet.WalletID = uuid.NewV4()
	}

	if wallets, err := interactor.storageDB.createWallets(newWallets); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating wallets on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return wallets, nil
	}
}

// updateWallet ...
func (interactor *interactor) updateWallet(updWallet *wallet) (*wallet, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "updateWallet"})
	log.Infof("updating wallet %s of user %s", updWallet.UserID.String(), updWallet.UserID.String())
	if wallet, err := interactor.storageDB.updateWallet(updWallet); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error updating wallet on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return wallet, nil
	}
}

// deleteWallet ...
func (interactor *interactor) deleteWallet(userID uuid.UUID, walletID uuid.UUID) *goerror.ErrorData {
	log.WithFields(map[string]interface{}{"method": "deleteWallet"})
	log.Infof("deleting wallet %s of user %s", walletID.String(), userID.String())
	if err := interactor.storageDB.deleteWallet(userID, walletID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting wallet on storage database %s", err).ToErrorData(err)
		return err
	}
	return nil
}

// getImages ...
func (interactor *interactor) getImages(userID uuid.UUID) ([]*image, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getImages"})
	log.Infof("getting images of user %s", userID.String())
	if images, err := interactor.storageDB.getImages(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting images on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return images, nil
	}
}

// getImage ...
func (interactor *interactor) getImage(userID uuid.UUID, imageID uuid.UUID) (*image, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getImage"})
	log.Infof("getting image %s of user %s", imageID.String(), userID.String())

	if image, err := interactor.storageDB.getImage(userID, imageID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting image on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return image, nil
	}
}

// getImageRaw ...
func (interactor *interactor) getImageRaw(userID uuid.UUID, imageID uuid.UUID) ([]byte, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getImage"})
	log.Infof("getting rawImage %s of user %s", imageID.String(), userID.String())

	if interactor.config.Dropbox.Enabled {
		path := fmt.Sprintf("/users/%s/images/%s", userID, imageID)
		if rawImage, err := interactor.storageDropbox.download(path); err != nil {
			log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
				Errorf("error getting rawImage on storage dropbox %s", err).ToErrorData(err)
			return nil, err
		} else {
			return rawImage, nil
		}
	} else {
		if image, err := interactor.storageDB.getImage(userID, imageID); err != nil {
			log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
				Errorf("error getting image on storage database %s", err).ToErrorData(err)
			return nil, err
		} else {
			return image.RawImage, nil
		}
	}
}

// createImage ...
func (interactor *interactor) createImage(newImage *image) (*image, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "createImage"})

	log.Info("creating image")
	newImage.ImageID = uuid.NewV4()

	if image, err := interactor.storageDB.createImage(newImage); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating image on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		if interactor.config.Dropbox.Enabled {
			path := fmt.Sprintf("/users/%s/images/%s", newImage.UserID.String(), newImage.ImageID.String())
			if err := interactor.storageDropbox.upload(path, newImage.RawImage); err != nil {
				log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
					Errorf("error creating image on storage dropbox %s", err).ToErrorData(err)
				return nil, err
			}
		}

		return image, nil
	}
}

// updateImage ...
func (interactor *interactor) updateImage(updImage *image) (*image, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "updateImage"})
	log.Infof("updating image %s of user %s", updImage.UserID.String(), updImage.UserID.String())
	if image, err := interactor.storageDB.updateImage(updImage); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error updating image on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		if interactor.config.Dropbox.Enabled {
			path := fmt.Sprintf("/users/%s/images/%s", updImage.UserID.String(), updImage.ImageID.String())
			if err := interactor.storageDropbox.upload(path, updImage.RawImage); err != nil {
				log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
					Errorf("error updating image on storage dropbox %s", err).ToErrorData(err)
				return nil, err
			}
			return image, nil
		}
	}
	return nil, nil
}

// deleteImage ...
func (interactor *interactor) deleteImage(userID uuid.UUID, imageID uuid.UUID) *goerror.ErrorData {
	log.WithFields(map[string]interface{}{"method": "deleteImage"})
	log.Infof("deleting image %s of user %s", imageID.String(), userID.String())
	if err := interactor.storageDB.deleteImage(userID, imageID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting image on storage database %s", err).ToErrorData(err)
		return err
	} else {
		if interactor.config.Dropbox.Enabled {
			path := fmt.Sprintf("/users/%s/images/%s", userID, imageID)
			if err := interactor.storageDropbox.delete(path); err != nil {
				log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
					Errorf("error deleting image on storage dropbox %s", err).ToErrorData(err)
				return err
			}
		}
	}
	return nil
}

// getCategories ...
func (interactor *interactor) getCategories(userID uuid.UUID) ([]*category, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getCategories"})
	log.Infof("getting categories of user %s", userID.String())
	if categories, err := interactor.storageDB.getCategories(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting categories on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return categories, nil
	}
}

// getCategory ...
func (interactor *interactor) getCategory(userID uuid.UUID, categoryID uuid.UUID) (*category, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getCategory"})
	log.Infof("getting category %s of user %s", categoryID.String(), userID.String())
	if category, err := interactor.storageDB.getCategory(userID, categoryID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting category on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return category, nil
	}
}

// createCategories ...
func (interactor *interactor) createCategories(newCategories []*category) ([]*category, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "createCategories"})

	log.Info("creating categories")
	for _, category := range newCategories {
		category.CategoryID = uuid.NewV4()
	}

	if categories, err := interactor.storageDB.createCategories(newCategories); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating categories on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return categories, nil
	}
}

// updateCategory ...
func (interactor *interactor) updateCategory(updCategory *category) (*category, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "updateCategory"})
	log.Infof("updating category %s of user %s", updCategory.UserID.String(), updCategory.UserID.String())
	if category, err := interactor.storageDB.updateCategory(updCategory); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error updating category on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return category, nil
	}
}

// deleteCategory ...
func (interactor *interactor) deleteCategory(userID uuid.UUID, categoryID uuid.UUID) *goerror.ErrorData {
	log.WithFields(map[string]interface{}{"method": "deleteCategory"})
	log.Infof("deleting category %s of user %s", categoryID.String(), userID.String())
	if err := interactor.storageDB.deleteCategory(userID, categoryID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting category on storage database %s", err).ToErrorData(err)
		return err
	}
	return nil
}

// getTransaction ...
func (interactor *interactor) getTransactions(userID uuid.UUID) ([]*transaction, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getTransactions"})
	log.Infof("getting transactions of user %s", userID.String())
	if transaction, err := interactor.storageDB.getTransactions(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting transactions on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// getTransaction ...
func (interactor *interactor) getTransaction(userID uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID) (*transaction, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "getTransaction"})
	log.Infof("getting transaction %s of user %s on wallet %s", transactionID.String(), userID.String(), walletID.String())
	if transaction, err := interactor.storageDB.getTransaction(userID, walletID, transactionID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting transaction on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// createTransactions ...
func (interactor *interactor) createTransactions(newTransactions []*transaction) ([]*transaction, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "createTransactions"})

	log.Info("creating transactions")
	for _, transaction := range newTransactions {
		transaction.TransactionID = uuid.NewV4()
	}

	if transactions, err := interactor.storageDB.createTransactions(newTransactions); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating transactions on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return transactions, nil
	}
}

// updateTransaction ...
func (interactor *interactor) updateTransaction(updTransaction *transaction) (*transaction, *goerror.ErrorData) {
	log.WithFields(map[string]interface{}{"method": "updateTransaction"})
	log.Infof("updating transaction %s of user %s", updTransaction.UserID.String(), updTransaction.UserID.String())
	if transaction, err := interactor.storageDB.updateTransaction(updTransaction); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error updating transaction on storage database %s", err).ToErrorData(err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// deleteTransaction ...
func (interactor *interactor) deleteTransaction(userID uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID) *goerror.ErrorData {
	log.WithFields(map[string]interface{}{"method": "deleteTransaction"})
	log.Infof("deleting transaction %s of user %s", transactionID.String(), userID.String())
	if err := interactor.storageDB.deleteTransaction(userID, walletID, transactionID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting transaction on storage database %s", err).ToErrorData(err)
		return err
	}
	return nil
}