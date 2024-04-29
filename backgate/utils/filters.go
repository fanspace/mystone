package utils

import (
	"backgate/validations"
	"strings"
)

func Formatphone(phone string) string {
	if phone == "" {
		return phone
	}
	if len(phone) < 8 {
		return phone
	}
	return phone[0:3] + "****" + phone[7:]
}

func FormatEmailAddress(email string) string {
	if email == "" {
		return email
	}
	if strings.Index(email, "@") < 0 {
		return email
	}
	mailarr := strings.Split(email, "@")
	mprefix := mailarr[0]
	if iscn := validations.IsChineseChar(mprefix); iscn {
		return ""
	}
	if len(mprefix) > 3 {
		return mprefix[0:3] + "***@" + mailarr[1]
	} else {
		return email
	}

}

func FormatMobileNo(mobile string) string {
	if mobile == "" {
		return mobile
	}
	if len(mobile) != 11 {
		return ""
	}
	return mobile[0:3] + "****" + mobile[7:11]
}

func FormatName(name string) string {
	if name == "" {
		return name
	}
	if len(name) < 2 {
		return name
	}
	return name[0:1] + "***"
}

func FormatUsername(username string) string {
	if username == "" {
		return username
	}

	if len(username) < 3 {
		return username
	}

	if len(username) < 6 {
		return username[0:1] + "***" + username[len(username)-1:]
	}
	return username[0:2] + "***" + username[len(username)-2:]
}
