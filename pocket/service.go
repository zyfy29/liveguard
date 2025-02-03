package pocket

func (c *Client) ValidateToken() (string, error) {
	content, err := formatter[UserInfoItem]{c}.doRestWithResult("", "https://pocketapi.48.cn/user/api/v1/user/info/reload", map[string]string{})
	if err != nil {
		return "", err
	}
	return content.Nickname, nil
}

func (c *Client) GetHistoryLives(memberID string, nextTime int64) ([]LiveItem, int64, error) {
	content, err := formatter[LiveListContent]{c}.doRestWithResult("", "https://pocketapi.48.cn/im/api/v1/chatroom/msg/list/aim/type", newLiveListReq(memberID, nextTime))
	if err != nil {
		return nil, 0, err
	}
	var lives []LiveItem
	for _, item := range content.Message {
		lives = append(lives, item.FormatToLiveItem())
	}
	return lives, content.NextTime, nil
}

func (c *Client) GetLiveInfo(liveId string) (LiveInfoContent, error) {
	content, err := formatter[LiveInfoContent]{c}.doRestWithResult("", "https://pocketapi.48.cn/live/api/v1/live/getLiveOne", map[string]string{
		"liveId": liveId,
	})
	if err != nil {
		return LiveInfoContent{}, err
	}
	return content, nil
}
