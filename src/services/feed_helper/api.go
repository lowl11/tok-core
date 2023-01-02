package feed_helper

import (
	"github.com/lowl11/lazy-collection/array"
	"sync"
	"tok-core/src/data/entities"
	"tok-core/src/repositories/post_comment_repository"
	"tok-core/src/repositories/post_like_repository"
)

func SetCommentRepository(commentRepository *post_comment_repository.Repository) {
	commentRepo = commentRepository
}

func SetLikeRepository(likeRepository *post_like_repository.Repository) {
	likeRepo = likeRepository
}

/*
	LikesAndComments возвращает список лайков и комментариев по заданным кодам постов
*/
func LikesAndComments(postCodes []string) (*array.Array[entities.PostLikeGetList], *array.Array[entities.PostCommentGetList]) {
	// запишем сходу слайсы т.к. и получаем мы их так же разом из репозитория
	likeChannel := make(chan []entities.PostLikeGetList)
	commentChannel := make(chan []entities.PostCommentGetList)

	// счетчик ставим на 2: лайки и комментарии
	wg := sync.WaitGroup{}
	wg.Add(2)

	// сходу запускаем wait, после чьего выполнения закроются каналы
	go func() {
		wg.Wait()
		close(likeChannel)
		close(commentChannel)
	}()

	// непосредственно методы получения
	go getLikes(&wg, likeChannel, postCodes)
	go getComments(&wg, commentChannel, postCodes)

	// записываем листы после того как забрали данные
	likeList := <-likeChannel
	commentList := <-commentChannel

	var likeArray *array.Array[entities.PostLikeGetList]
	if likeList != nil && len(likeList) > 0 {
		likeArray = array.NewWithList[entities.PostLikeGetList](likeList...)
	}

	var commentArray *array.Array[entities.PostCommentGetList]
	if commentList != nil && len(commentList) > 0 {
		commentArray = array.NewWithList[entities.PostCommentGetList](commentList...)
	}

	return likeArray, commentArray
}
