package param

type RegisterParams struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}
type LoginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type PostsPageParams struct {
	CommunityID int64  `form:"community_id"`
	Page        int    `form:"page"`
	Size        int    `form:"size"`
	Order       string `form:"order"`
}
