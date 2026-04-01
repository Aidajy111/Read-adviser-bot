package files

import (
	"encoding/gob"
	"errors"
	"math/rand/v2"
	"os"
	"path/filepath"

	"github.com/Aidajy111/Read-adviser-bot/lib/e"
	"github.com/Aidajy111/Read-adviser-bot/storage"
)

type Storage struct {
	basePath string
}

const (
	defaultPerm        = 0774
	errSaveMsg         = "can t sav pagee"
	errRandomPickMasg  = "can t random page"
	errCantDecodePage  = "can t decode page"
	errRemovePageMsg   = "can t remove page"
	errIsExistsPageMsg = "can t check is exists page"
)

var ErrSavesPages = errors.New("no saved pages")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return e.Wrap(errSaveMsg, err)
	}

	fName, err := fileName(page)
	if err != nil {
		return e.Wrap(errSaveMsg, err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap(errSaveMsg, err)
	}

	defer func() {
		_ = file.Close()
	}()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap(errSaveMsg, err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.Wrap(errRandomPickMasg, err)
	}

	if len(files) == 0 {
		return nil, ErrSavesPages
	}

	n := rand.IntN(len(files))

	file := files[n]

	// open docede
	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fileName, err := fileName(page)
	if err != nil {
		return e.Wrap(errRemovePageMsg, err)
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return e.Wrap(errRemovePageMsg, err)
	}

	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fileName, err := fileName(page)
	if err != nil {
		return false, e.Wrap(errIsExistsPageMsg, err)
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(errIsExistsPageMsg, err)
	default:
		return true, nil
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap(errCantDecodePage, err)
	}

	defer func() {
		_ = f.Close()
	}()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap(errCantDecodePage, err)
	}
	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
