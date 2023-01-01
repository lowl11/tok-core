package feed_helper

import (
	"github.com/lowl11/lazylog/layers"
	"sync"
	"tok-core/src/data/entities"
	"tok-core/src/definition"
	"tok-core/src/repositories/post_comment_repository"
	"tok-core/src/repositories/post_like_repository"
)

var (
	commentRepo *post_comment_repository.Repository
	likeRepo    *post_like_repository.Repository
)

func getComments(wg *sync.WaitGroup, commentChannel chan []entities.PostCommentGetList, postCodes []string) {
	defer wg.Done()

	logger := definition.Logger

	// получаем комменты постов
	commentList, err := commentRepo.GetByList(postCodes)
	if err != nil {
		logger.Error(err, "Get posts comment counter error", layers.Mongo)
	}

	commentChannel <- commentList
}

func getLikes(wg *sync.WaitGroup, likeChannel chan []entities.PostLikeGetList, postCodes []string) {
	defer wg.Done()

	logger := definition.Logger

	// получаем лайки постов
	likeList, err := likeRepo.GetByList(postCodes)
	if err != nil {
		logger.Error(err, "Get posts like counter error", layers.Mongo)
		return
	}

	likeChannel <- likeList
}
