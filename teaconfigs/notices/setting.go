package notices

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

// 通知设置
type NoticeSetting struct {
	Levels map[NoticeLevel]*NoticeLevelConfig `yaml:"levels" json:"levels"`
	Medias []*NoticeMediaConfig               `yaml:"medias" json:"medias"`
}

// 取得当前的配置
func SharedNoticeSetting() *NoticeSetting {
	filename := "notice.conf"
	file := files.NewFile(Tea.ConfigFile(filename))
	config := &NoticeSetting{
		Levels: map[NoticeLevel]*NoticeLevelConfig{},
		Medias: []*NoticeMediaConfig{},
	}
	if !file.Exists() {
		return config
	}

	reader, err := file.Reader()
	if err != nil {
		logs.Error(err)
		return config
	}
	defer reader.Close()
	err = reader.ReadYAML(config)
	if err != nil {
		logs.Error(err)
		return config
	}
	return config
}

// 保存配置
func (this *NoticeSetting) Save() error {
	filename := "notice.conf"
	writer, err := files.NewWriter(Tea.ConfigFile(filename))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 查找级别配置
func (this *NoticeSetting) LevelConfig(level NoticeLevel) *NoticeLevelConfig {
	config, found := this.Levels[level]
	if found {
		return config
	}
	config = &NoticeLevelConfig{
		ShouldNotify: true,
	}
	this.Levels[level] = config
	return config
}

// 添加媒介配置
func (this *NoticeSetting) AddMedia(mediaConfig *NoticeMediaConfig) {
	this.Medias = append(this.Medias, mediaConfig)
}

// 删除媒介
func (this *NoticeSetting) RemoveMedia(mediaId string) {
	medias := []*NoticeMediaConfig{}
	for _, m := range this.Medias {
		if m.Id == mediaId {
			continue
		}
		medias = append(medias, m)
	}
	this.Medias = medias

	// 移除关联的接收人
	for _, l := range this.Levels {
		l.RemoveMediaReceivers(mediaId)
	}
}

// 查找媒介
func (this *NoticeSetting) FindMedia(mediaId string) *NoticeMediaConfig {
	for _, m := range this.Medias {
		if m.Id == mediaId {
			m.Validate()
			return m
		}
	}
	return nil
}

// 查找接收人
func (this *NoticeSetting) FindReceiver(receiverId string) (level NoticeLevel, receiver *NoticeReceiver) {
	for levelCode, levelConfig := range this.Levels {
		receiver := levelConfig.FindReceiver(receiverId)
		if receiver != nil {
			return levelCode, receiver
		}
	}
	return 0, nil
}

// 发送通知
func (this *NoticeSetting) Notify(level NoticeLevel, message string, counter func(receiverId string, minutes int) int) (receiverIds []string) {
	config, found := this.Levels[level]
	if !found {
		return
	}
	for _, r := range config.Receivers {
		if !r.On {
			continue
		}
		media := this.FindMedia(r.MediaId)
		if !media.On {
			continue
		}
		mediaType := FindNoticeMediaType(media.Type)
		if mediaType == nil {
			continue
		}
		raw, err := media.Raw()
		if err != nil {
			logs.Error(err)
			continue
		}
		if !media.ShouldNotify(counter(r.Id, media.RateMinutes)) {
			continue
		}
		receiverIds = append(receiverIds, r.Id)
		go func(raw NoticeMediaInterface, mediaType maps.Map, user string) {
			body := message
			if types.Bool(mediaType["supportsHTML"]) {
				body = strings.Replace(body, "\n", "<br/>", -1)
			}
			_, err := raw.Send(user, "[TeaWeb]["+FindNoticeLevelName(level)+"]有新的通知", body)
			if err != nil {
				logs.Error(err)
			}
		}(raw, mediaType, r.User)
	}
	return
}