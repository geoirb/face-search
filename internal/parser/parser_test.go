package parser

import (
	"testing"

	service "github.com/geoirb/face-search/internal/face-search"
	"github.com/stretchr/testify/assert"
)

var (
	profileLayout = `<div class="card-vk01-header">([^<]*)<\/div><div class="card-vk01-score">Совпадение: <span class="score-label">([0-9]{1,2}[.]?[0-9]{1,2}%)<\/span><\/div><div class="[^<]*">[^<]*<\/div><div class="card-vk01-geo">[^<]*<\/div><div class="btn-vk01-container"><a href="(https:\/\/vk.com\/[^"]*)" target="_blank" class="btn btn-primary btn-vk01">Профиль<\/a><a href="#" data-target="#modalIMG" data-toggle="modal" class="btn btn-primary btn-vk01" data-imgsrc="[^"]*" data-imghref="(https:\/\/vk.com\/[^"]*)">Фото<\/a>`
	payload       = []byte(`<div class="card card-vk01  border border-primary"><div class="row no-gutters"><div class="class=" card-vk01-fixed"=""><a href="https://vk.com/id223636747?z=photo223636747_456251236%2Fphotos223636747" target="_blank"><img src="//i.search4faces.com/faces/vk01/07/69/65/769656823/0.jpg" class="card-img-vk01" alt="Сергей Сыроежкин "></a></div><div class="col"><div class="card-body card-vk01-body"><div class="card-vk01-header">Сергей Сыроежкин </div><div class="card-vk01-score">Совпадение: <span class="score-label">52.69%</span></div><div class="card-vk01-age">35 лет | 8.8.1985</div><div class="card-vk01-geo">Россия, Москва</div><div class="btn-vk01-container"><a href="https://vk.com/id223636747" target="_blank" class="btn btn-primary btn-vk01">Профиль</a><a href="#" data-target="#modalIMG" data-toggle="modal" class="btn btn-primary btn-vk01" data-imgsrc="https://sun9-10.userapi.com/c831409/v831409120/1098c1/4wcKmkpOgFE.jpg" data-imghref="https://vk.com/id223636747?z=photo223636747_456251236%2Fphotos223636747">Фото</a></div></div></div></div></div><div class="card card-vk01  border border-primary"><div class="row no-gutters"><div class="class=" card-vk01-fixed"=""><a href="https://vk.com/id3746498?z=photo3746498_239278744%2Fphotos3746498" target="_blank"><img src="//i.search4faces.com/faces/vk01/02/31/15/231151356/1.jpg" class="card-img-vk01" alt="Ірина Тіпстеровіч "></a></div><div class="col"><div class="card-body card-vk01-body"><div class="card-vk01-header">Ірина Тіпстеровіч </div><div class="card-vk01-score">Совпадение: <span class="score-label">52.68%</span></div><div class="card-vk01-age">33 года | 11.8.1987</div><div class="card-vk01-geo">Украина, Киев</div><div class="btn-vk01-container"><a href="https://vk.com/id3746498" target="_blank" class="btn btn-primary btn-vk01">Профиль</a><a href="#" data-target="#modalIMG" data-toggle="modal" class="btn btn-primary btn-vk01" data-imgsrc="https://sun9-63.userapi.com/c10100/u3746498/101355907/w_9df3d059.jpg" data-imghref="https://vk.com/id3746498?z=photo3746498_239278744%2Fphotos3746498">Фото</a></div></div></div></div></div><div class="card card-vk01  border border-primary"><div class="row no-gutters"><div class="class=" card-vk01-fixed"=""><a href="https://vk.com/id132286728?z=photo132286728_338898962%2Fphotos132286728" target="_blank"><img src="//i.search4faces.com/faces/vk01/04/47/40/447404568/0.jpg" class="card-img-vk01" alt="Николай Дубов "></a></div><div class="col"><div class="card-body card-vk01-body"><div class="card-vk01-header">Николай Дубов </div><div class="card-vk01-score">Совпадение: <span class="score-label">52.30%</span></div><div class="card-vk01-age"></div><div class="card-vk01-geo">Россия, Челябинск</div><div class="btn-vk01-container"><a href="https://vk.com/id132286728" target="_blank" class="btn btn-primary btn-vk01">Профиль</a><a href="#" data-target="#modalIMG" data-toggle="modal" class="btn btn-primary btn-vk01" data-imgsrc="https://sun9-22.userapi.com/c623826/v623826728/125f/m095THSR6YI.jpg" data-imghref="https://vk.com/id132286728?z=photo132286728_338898962%2Fphotos132286728">Фото</a></div></div></div></div></div><div class="card card-vk01  border border-primary"><div class="row no-gutters"><div class="class=" card-vk01-fixed"=""><a href="https://vk.com/id552573271?z=photo552573271_457243038%2Fphotos552573271" target="_blank"><img src="//i.search4faces.com/faces/vk01/13/15/80/1315801925/0.jpg" class="card-img-vk01" alt="有毒的寄生虫 我 "></a></div>`)

	expectedProfile = []service.Profile{
		{
			FullName:    "Сергей Сыроежкин ",
			LinkProfile: "https://vk.com/id223636747",
			LinkPhoto:   "https://vk.com/id223636747?z=photo223636747_456251236%2Fphotos223636747",
			Confidence:  "52.69%",
		},
		{
			FullName:    "Ірина Тіпстеровіч ",
			LinkProfile: "https://vk.com/id3746498",
			LinkPhoto:   "https://vk.com/id3746498?z=photo3746498_239278744%2Fphotos3746498",
			Confidence:  "52.68%",
		},
		{
			FullName:    "Николай Дубов ",
			LinkProfile: "https://vk.com/id132286728",
			LinkPhoto:   "https://vk.com/id132286728?z=photo132286728_338898962%2Fphotos132286728",
			Confidence:  "52.30%",
		},
	}
)

func TestParse(t *testing.T) {
	p, err := New(profileLayout)
	assert.NoError(t, err)
	actualProfiles := p.GetProfileList(payload)
	assert.Equal(t, expectedProfile, actualProfiles)
}
