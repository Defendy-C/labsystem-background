package admin

import (
	"go.uber.org/zap"
	"labsystem/model"
	"labsystem/util"
	"labsystem/util/logger"
	"labsystem/util/rsa"
	"time"
)

type loginReq struct {
	AdminNick string `json:"user_name"`
	Password  string `json:"password"`
	Key       string `json:"key"`
	VCode     int    `json:"v_code"`
}

func (req *loginReq) Valid() bool {
	if err := util.StringFormatVerify(req.AdminNick, model.RegExpUserName); err != nil {
		logger.Log.Warn("admin nick verify error:", zap.Error(err), zap.String("adminNick", req.AdminNick))
		return false
	}

	return true
}

type InfoResp struct {
	Id     uint          `json:"id"`
	Name   string        `json:"name"`
	Powers []*PowerOwner `json:"powers"`
}

type PowerOwner struct {
	Name  string      `json:"name"`
	Power model.Power `json:"power"`
	Own   bool        `json:"own"`
}

type ListReq struct {
	CreatedBy uint   `json:"created_by"`
	Page      uint   `json:"page"`
	PageSize  uint   `json:"page_size"`
}

type Item struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	Power     []*PowerOwner `json:"power"`
	CreatedBy string        `json:"created_by"`
	CreatedAt time.Time     `json:"created_at"`
}

type ListResp struct {
	List       []*Item `json:"list"`
	TotalPage  uint    `json:"total_page"`
	TotalCount uint    `json:"total_count"`
}

type CreateAdminReq struct {
	Name string `json:"name"`
	Password string `json:"password"`
	Power int `json:"power"`
}

func (req *CreateAdminReq)Valid() bool {
	var err error
	if _, err = model.IntToPower(req.Power); err != nil {
		return false
	}
	if err = util.StringFormatVerify(req.Name, model.RegExpUserName); err != nil {
		return false
	}
	var pwd string
	if pwd, err = rsa.Decrypt(req.Password); err != nil {
		return false
	}
	if err := util.StringFormatVerify(pwd, model.RegExpPassword); err != nil {
		return false
	}

	return true
}

type UpdateAdminReq struct {
	Name string `json:"name"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

func (req *UpdateAdminReq)Valid() bool {
	var err error
	if err = util.StringFormatVerify(req.Name, model.RegExpUserName); err != nil {
		return false
	}
	var newPwd, oldPwd string
	if newPwd != "" {
		if newPwd, err = rsa.Decrypt(req.NewPassword); err != nil {
			return false
		}
		if oldPwd, err = rsa.Decrypt(req.OldPassword); err != nil {
			return false
		}
		if err := util.StringFormatVerify(newPwd, model.RegExpPassword); err != nil {
			return false
		}
		if err := util.StringFormatVerify(oldPwd, model.RegExpPassword); err != nil {
			return false
		}
		if newPwd == oldPwd {
			return false
		}
	}

	return true
}

type DeleteAdminReq struct {
	ID uint `json:"id"`
}

type CreateTeacherReq struct {
	UserNo   string `json:"user_no"`
	RealName string `json:"real_name"`
	Password string `json:"password"`
	Class    string `json:"class"`
	FileName string `json:"file_name"`
}

func (req *CreateTeacherReq) Valid() bool {
	if err := util.StringFormatVerify(req.UserNo, model.RegExpUserNo); err != nil {
		return false
	}
	if err := util.StringFormatVerify(req.RealName, model.ReqExpRealName); err != nil {
		return false
	}
	pwd, err := rsa.Decrypt(req.Password)
	if err != nil {
		return false
	}
	if err := util.StringFormatVerify(pwd, model.RegExpPassword); err != nil {
		return false
	}

	return true
}

type CreateClassReq struct {
	ClassNo string `json:"class_no"`
}

func (req *CreateClassReq) Valid() bool {
	err := util.StringFormatVerify(req.ClassNo, model.RegExpClass)
	if err != nil {
		return false
	}

	return true
}

type UserListReq struct {
	Page uint `json:"page"`
	PageSize uint `json:"page_size"`
}

type UserItem struct {
	ID uint `json:"id"`
	UserNo string `json:"user_no"`
	RealName string `json:"real_name"`
	Status model.UserStatus `json:"status"`
	Class string `json:"class"`
	ProfileUrl string `json:"profile_url"`
	CreatedBy string `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type UserListResp struct {
	List []*UserItem `json:"list"`
	TotalPage uint `json:"total_page"`
	TotalCount uint `json:"total_count"`
}

type ClassListReq struct {
	Page uint `json:"page"`
	PageSize uint `json:"page_size"`
}

type ClassItem struct {
	ID uint `json:"id"`
	ClassNo string `json:"class_no"`
	CreatedAt time.Time `json:"created_at"`
}

type ClassListResp struct {
	List []*ClassItem `json:"list"`
	TotalPage uint `json:"total_page"`
	TotalCount uint `json:"total_count"`
}

type DeleteUsersReq struct {
	Ids []uint `json:"ids"`
}