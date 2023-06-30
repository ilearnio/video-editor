package shotcutProjectBuilder

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"videoeditor/src/models"
)

var _ = Describe("Shotcut Project Builder", Ordered, func() {
	It("should build a simple MLT", func() {
		profile := Profile{
			Description: "PAL 4:3 DV or DVD",
		}
		playlist := Playlist{
			Id: "playlist0",
			Children: []XMLConvertible{
				Property{"shotcut:name", "1"},
				Property{"shotcut:audio", "2"},
			},
		}
		tractor := Tractor{
			Id:       "tractor0",
			Title:    "",
			In:       "",
			Out:      "",
			Children: nil,
		}
		mlt := MLT{
			Version:    "7.17.0",
			LC_NUMERIC: "C",
			Title:      "",
			Producer:   "",
			Children:   []XMLConvertible{profile, playlist, tractor},
		}
		xml := NodeToXML(mlt.ToXMLNode(), 0)

		Expect(xml).To(Equal(`<mlt LC_NUMERIC="C" version="7.17.0">
  <profile description="PAL 4:3 DV or DVD"/>
  <playlist id="playlist0">
    <property name="shotcut:name">1</property>
    <property name="shotcut:audio">2</property>
  </playlist>
  <tractor id="tractor0"/>
</mlt>`))
	})

	It("should build entire .mlt file", func() {
		backgroundImagePath := "slides/bg.jpg"
		introImagePath := "/Users/admin/Documents/Videos/assets/intro/motivation-daily-intro.jpg"
		outroImagePath := "/Users/admin/Documents/Videos/assets/outro/motivations-daily-outro.jpg"
		outroOverlayImagePath := "/Users/admin/Documents/Videos/assets/outro/motivations-daily-outro-2.jpg"

		quoteAudio := AssetAudio{
			Path:     "audio/01 - In school we lea 20 (enhanced).wav",
			Duration: 23440 * time.Millisecond,
		}
		quoteAudio2 := AssetAudio{
			Path:     "audio/02 - Winners are not 2 (enhanced).wav",
			Duration: 14520 * time.Millisecond,
		}
		quoteAudio3 := AssetAudio{
			Path:     "audio/02 - Winners are not 2 (enhanced).wav",
			Duration: 14520 * time.Millisecond,
		}
		quoteAudios := []AssetAudio{quoteAudio, quoteAudio2, quoteAudio3}

		quotes := []*models.VideoQuote{
			{Content: "quote text 1", IsHtmlEnabled: false},
			{Content: "quote text 2", IsHtmlEnabled: false},
			{Content: "quote text 3", IsHtmlEnabled: false},
		}

		bgMusicAudio := AssetAudio{
			Path:     "/Users/admin/Documents/Videos/sources/audio/1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3",
			Duration: 30000 * time.Millisecond,
		}

		xml, err := BuildQuotesProject(BuildQuotesProjectParams{
			BackgroundImagePath:   backgroundImagePath,
			IntroImagePath:        introImagePath,
			OutroImagePath:        outroImagePath,
			OutroOverlayImagePath: outroOverlayImagePath,
			HeadingIsHTML:         true,
			HeadingContent:        "Title <br> test",
			QuoteAudios:           quoteAudios,
			Quotes:                quotes,
			BgMusicAudio:          bgMusicAudio,
		})

		expectedXML := `<?xml version="1.0" standalone="no"?>
<mlt LC_NUMERIC="C" version="7.17.0" title="Shotcut version 23.05.14" producer="main_bin">
  <profile description="PAL 4:3 DV or DVD" width="1920" height="1080" progressive="1" sample_aspect_num="1" sample_aspect_den="1" display_aspect_num="16" display_aspect_den="9" frame_rate_num="25" frame_rate_den="1" colorspace="709"/>
  <playlist id="main_bin">
    <property name="xml_retain">1</property>
  </playlist>
  <producer id="producer_background_color" in="00:00:00.000" out="00:01:07.480">
    <property name="length">00:01:07.480</property>
    <property name="eof">pause</property>
    <property name="resource">0</property>
    <property name="aspect_ratio">1</property>
    <property name="mlt_service">color</property>
    <property name="mlt_image_format">rgba</property>
    <property name="set.test_audio">0</property>
  </producer>
  <playlist id="playlist_background_color">
    <property name="shotcut:name">Background Color</property>
    <entry producer="producer_background_color" in="00:00:00.000" out="00:01:07.480"/>
  </playlist>
  <producer id="bg_image_producer" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">bg.jpg</property>
    <property name="ttl">1</property>
    <property name="aspect_ratio">1</property>
    <property name="progressive">1</property>
    <property name="seekable">1</property>
    <property name="format">1</property>
    <property name="mlt_service">qimage</property>
    <property name="creation_time">2023-05-23T13:36:16</property>
    <property name="shotcut:caption">bg.jpg</property>
    <property name="xml">was here</property>
  </producer>
  <playlist id="playlist0">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Background Image</property>
    <entry producer="bg_image_producer" in="00:00:00.000" out="00:01:07.480"/>
  </playlist>
  <chain id="chain0" out="00:00:23.440">
    <property name="length">00:00:23.440</property>
    <property name="eof">pause</property>
    <property name="resource">quotes/01 - In school we lea 20 (enhanced).wav</property>
    <property name="mlt_service">avformat-novalidate</property>
    <property name="seekable">1</property>
    <property name="audio_index">0</property>
    <property name="video_index">-1</property>
    <property name="mute_on_pause">0</property>
    <property name="video_delay">0</property>
    <property name="shotcut:caption">01 - In school we lea 20 (enhanced).wav</property>
    <property name="shotcut:producer">avformat-novalidate</property>
    <property name="xml">was here</property>
  </chain>
  <chain id="chain1" out="00:00:14.520">
    <property name="length">00:00:14.520</property>
    <property name="eof">pause</property>
    <property name="resource">quotes/02 - Winners are not 2 (enhanced).wav</property>
    <property name="mlt_service">avformat-novalidate</property>
    <property name="seekable">1</property>
    <property name="audio_index">0</property>
    <property name="video_index">-1</property>
    <property name="mute_on_pause">0</property>
    <property name="video_delay">0</property>
    <property name="shotcut:caption">02 - Winners are not 2 (enhanced).wav</property>
    <property name="shotcut:producer">avformat-novalidate</property>
    <property name="xml">was here</property>
  </chain>
  <chain id="chain2" out="00:00:14.520">
    <property name="length">00:00:14.520</property>
    <property name="eof">pause</property>
    <property name="resource">quotes/02 - Winners are not 2 (enhanced).wav</property>
    <property name="mlt_service">avformat-novalidate</property>
    <property name="seekable">1</property>
    <property name="audio_index">0</property>
    <property name="video_index">-1</property>
    <property name="mute_on_pause">0</property>
    <property name="video_delay">0</property>
    <property name="shotcut:caption">02 - Winners are not 2 (enhanced).wav</property>
    <property name="shotcut:producer">avformat-novalidate</property>
    <property name="xml">was here</property>
  </chain>
  <playlist id="playlist1">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Quote Audio</property>
    <blank length="00:00:03.500"/>
    <entry producer="chain0" in="00:00:00.000" out="00:00:23.440"/>
    <blank length="00:00:02.000"/>
    <entry producer="chain1" out="00:00:14.520"/>
    <blank length="00:00:02.000"/>
    <entry producer="chain2" out="00:00:14.520"/>
  </playlist>
  <!-- First track for background music -->
  <chain id="channel_bgmusic_0" out="00:00:30.000">
    <property name="length">00:00:30.000</property>
    <property name="eof">pause</property>
    <property name="resource">1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3</property>
    <property name="mlt_service">avformat-novalidate</property>
    <property name="seekable">1</property>
    <property name="audio_index">0</property>
    <property name="video_index">-1</property>
    <property name="mute_on_pause">0</property>
    <property name="video_delay">0</property>
    <property name="shotcut:caption">1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3</property>
    <property name="shotcut:producer">avformat-novalidate</property>
    <property name="xml">was here</property>
    <filter id="filter_channel_bgmusic_0_fadeOut" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">00:00:18.000=0;00:00:20.000=-60</property>
      <property name="mlt_service">volume</property>
      <property name="shotcut:filter">fadeOutVolume</property>
      <property name="shotcut:animOut">00:00:02.000</property>
      <property name="disable">0</property>
    </filter>
    <filter id="filter_channel_bgmusic_0_fadeIn" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">00:00:00.000=-60;00:00:02.000=0</property>
      <property name="mlt_service">volume</property>
      <property name="shotcut:filter">fadeInVolume</property>
      <property name="shotcut:animIn">00:00:02.000</property>
    </filter>
    <filter id="filter_channel_bgmusic_0_volume" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">-11.1</property>
      <property name="mlt_service">volume</property>
    </filter>
  </chain>
  <chain id="channel_bgmusic_2" out="00:00:30.000">
    <property name="length">00:00:30.000</property>
    <property name="eof">pause</property>
    <property name="resource">1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3</property>
    <property name="mlt_service">avformat-novalidate</property>
    <property name="seekable">1</property>
    <property name="audio_index">0</property>
    <property name="video_index">-1</property>
    <property name="mute_on_pause">0</property>
    <property name="video_delay">0</property>
    <property name="shotcut:caption">1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3</property>
    <property name="shotcut:producer">avformat-novalidate</property>
    <property name="xml">was here</property>
    <filter id="filter_channel_bgmusic_2_fadeOut" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">00:00:18.000=0;00:00:20.000=-60</property>
      <property name="mlt_service">volume</property>
      <property name="shotcut:filter">fadeOutVolume</property>
      <property name="shotcut:animOut">00:00:02.000</property>
      <property name="disable">0</property>
    </filter>
    <filter id="filter_channel_bgmusic_2_fadeIn" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">00:00:00.000=-60;00:00:02.000=0</property>
      <property name="mlt_service">volume</property>
      <property name="shotcut:filter">fadeInVolume</property>
      <property name="shotcut:animIn">00:00:02.000</property>
    </filter>
    <filter id="filter_channel_bgmusic_2_volume" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">-11.1</property>
      <property name="mlt_service">volume</property>
    </filter>
  </chain>

  <playlist id="playlist_bgmusic_1">
    <property name="shotcut:audio">1</property>
    <property name="shotcut:name">Background Music #1</property>
    <entry producer="channel_bgmusic_0" in="00:00:00.000" out="00:00:20.000"/>
    <blank length="00:00:16.000"/>
    <entry producer="channel_bgmusic_2" in="00:00:00.000" out="00:00:20.000"/>
  </playlist>
  <producer id="producer_intro_image" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">motivation-daily-intro.jpg</property>
    <property name="ttl">1</property>
    <property name="aspect_ratio">1</property>
    <property name="progressive">1</property>
    <property name="seekable">1</property>
    <property name="format">1</property>
    <property name="mlt_service">qimage</property>
    <property name="creation_time">2023-05-23T00:26:22</property>
    <property name="xml">was here</property>
    <property name="shotcut:caption">motivation-daily-intro.jpg</property>
    <filter id="filter0" out="00:00:04.080">
      <property name="start">1</property>
      <property name="level">00:00:00.000=0;00:00:00.680=1</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeInBrightness</property>
      <property name="alpha">1</property>
      <property name="shotcut:animIn">00:00:00.720</property>
    </filter>
    <filter id="filter7" out="00:00:04.080">
      <property name="start">1</property>
      <property name="level">00:00:03.400=1;00:00:04.080=0</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeOutBrightness</property>
      <property name="alpha">1</property>
      <property name="shotcut:animOut">00:00:00.720</property>
    </filter>
  </producer>
  <producer id="producer_outro_image" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">motivations-daily-outro.jpg</property>
    <property name="ttl">1</property>
    <property name="aspect_ratio">1</property>
    <property name="progressive">1</property>
    <property name="seekable">1</property>
    <property name="format">1</property>
    <property name="mlt_service">qimage</property>
    <property name="creation_time">2023-05-23T00:42:12</property>
    <property name="shotcut:caption">motivations-daily-outro.jpg</property>
    <property name="xml">was here</property>
  </producer>
  <playlist id="playlist_outro_image">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Intro and Outro</property>
    <property name="shotcut:lock">0</property>
    <entry producer="producer_intro_image" in="00:00:00.000" out="00:00:03.500"/>
    <blank length="00:00:56.480"/>
    <entry producer="producer_outro_image" in="00:00:00.000" out="00:00:07.500"/>
  </playlist>
  <producer id="producer_post_outro_image" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">motivations-daily-outro-2.jpg</property>
    <property name="ttl">1</property>
    <property name="aspect_ratio">1</property>
    <property name="progressive">1</property>
    <property name="seekable">1</property>
    <property name="format">1</property>
    <property name="mlt_service">qimage</property>
    <property name="creation_time">2023-05-23T00:41:13</property>
    <property name="shotcut:caption">motivations-daily-outro-2.jpg</property>
    <property name="xml">was here</property>
    <filter id="filter8" out="00:00:06.000">
      <property name="start">1</property>
      <property name="level">1</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeInBrightness</property>
      <property name="alpha">00:00:00.000=0;00:00:00.640=1</property>
      <property name="shotcut:animIn">00:00:00.680</property>
    </filter>
  </producer>
  <playlist id="playlist_overlay">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Outro Overlay</property>
    <blank length="00:01:01.480"/>
    <entry producer="producer_post_outro_image" in="00:00:00.000" out="00:00:06.000"/>
  </playlist>
  <producer id="producer_quote_0" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">#00000000</property>
    <property name="aspect_ratio">1</property>
    <property name="mlt_service">color</property>
    <property name="mlt_image_format">rgba</property>
    <property name="shotcut:caption">transparent</property>
    <property name="shotcut:detail">transparent</property>
    <property name="ignore_points">0</property>
    <property name="xml">was here</property>
    <property name="seekable">1</property>
    <filter id="producer_filter_producer_quote_0" out="00:00:23.440">
      <property name="argument"/>
      <property name="geometry">902 373 904 546 1</property>
      <property name="family">Sans</property>
      <property name="size">48</property>
      <property name="weight">400</property>
      <property name="style">normal</property>
      <property name="fgcolour">0x000000ff</property>
      <property name="bgcolour">#00000000</property>
      <property name="olcolour">0x00000000</property>
      <property name="pad">0</property>
      <property name="halign">left</property>
      <property name="valign">top</property>
      <property name="outline">0</property>
      <property name="pixel_ratio">1</property>
      <property name="mlt_service">qtext</property>
      <property name="shotcut:filter">richText</property>
      <property name="html">&lt;!DOCTYPE HTML PUBLIC &quot;-//W3C//DTD HTML 4.0//EN&quot; &quot;http://www.w3.org/TR/REC-html40/strict.dtd&quot;&gt;
&lt;html&gt;&lt;head&gt;&lt;meta name=&quot;qrichtext&quot; content=&quot;1&quot; /&gt;&lt;meta charset=&quot;utf-8&quot; /&gt;&lt;style type=&quot;text/css&quot;&gt;
p, li { white-space: pre-wrap; }
hr { height: 1px; border-width: 0; }
li.unchecked::marker { content: &quot;\2610&quot;; }
li.checked::marker { content: &quot;\2612&quot;; }
&lt;/style&gt;&lt;/head&gt;&lt;body style=&quot; font-family:'.AppleSystemUIFont'; font-size:13pt; font-weight:400; font-style:normal;&quot;&gt;
&lt;p align=&quot;center&quot; style=&quot; margin-top:0px; margin-bottom:0px; margin-left:0px; margin-right:0px; -qt-block-indent:0; text-indent:0px; line-height:117%;&quot;&gt;&lt;span style=&quot; font-family:'Mate SC'; font-size:41pt; font-weight:600; color:#ffffff;&quot;&gt;“quote text 1”&lt;/p&gt;&lt;/body&gt;&lt;/html&gt;</property>
      <property name="shotcut:animIn">00:00:00.000</property>
      <property name="shotcut:animOut">00:00:00.000</property>
      <property name="disable">0</property>
    </filter>
    <filter id="filter_producer_quote_0_fadeIn" out="00:00:23.440">
      <property name="start">1</property>
      <property name="level">1</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeInBrightness</property>
      <property name="alpha">00:00:00.000=0;00:00:01.000=1</property>
      <property name="shotcut:animIn">00:00:01.000</property>
    </filter>
    <filter id="filter_producer_quote_0_fadeOut" out="00:00:23.440">
      <property name="start">1</property>
      <property name="level">1</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeOutBrightness</property>
      <property name="alpha">00:00:22.440=1;00:00:23.440=0</property>
      <property name="shotcut:animOut">00:00:01.000</property>
    </filter>
  </producer>
  <producer id="producer_quote_1" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">#00000000</property>
    <property name="aspect_ratio">1</property>
    <property name="mlt_service">color</property>
    <property name="mlt_image_format">rgba</property>
    <property name="shotcut:caption">transparent</property>
    <property name="shotcut:detail">transparent</property>
    <property name="ignore_points">0</property>
    <property name="xml">was here</property>
    <property name="seekable">1</property>
    <filter id="producer_filter_producer_quote_1" out="00:00:14.520">
      <property name="argument"/>
      <property name="geometry">902 373 904 546 1</property>
      <property name="family">Sans</property>
      <property name="size">48</property>
      <property name="weight">400</property>
      <property name="style">normal</property>
      <property name="fgcolour">0x000000ff</property>
      <property name="bgcolour">#00000000</property>
      <property name="olcolour">0x00000000</property>
      <property name="pad">0</property>
      <property name="halign">left</property>
      <property name="valign">top</property>
      <property name="outline">0</property>
      <property name="pixel_ratio">1</property>
      <property name="mlt_service">qtext</property>
      <property name="shotcut:filter">richText</property>
      <property name="html">&lt;!DOCTYPE HTML PUBLIC &quot;-//W3C//DTD HTML 4.0//EN&quot; &quot;http://www.w3.org/TR/REC-html40/strict.dtd&quot;&gt;
&lt;html&gt;&lt;head&gt;&lt;meta name=&quot;qrichtext&quot; content=&quot;1&quot; /&gt;&lt;meta charset=&quot;utf-8&quot; /&gt;&lt;style type=&quot;text/css&quot;&gt;
p, li { white-space: pre-wrap; }
hr { height: 1px; border-width: 0; }
li.unchecked::marker { content: &quot;\2610&quot;; }
li.checked::marker { content: &quot;\2612&quot;; }
&lt;/style&gt;&lt;/head&gt;&lt;body style=&quot; font-family:'.AppleSystemUIFont'; font-size:13pt; font-weight:400; font-style:normal;&quot;&gt;
&lt;p align=&quot;center&quot; style=&quot; margin-top:0px; margin-bottom:0px; margin-left:0px; margin-right:0px; -qt-block-indent:0; text-indent:0px; line-height:117%;&quot;&gt;&lt;span style=&quot; font-family:'Mate SC'; font-size:41pt; font-weight:600; color:#ffffff;&quot;&gt;“quote text 2”&lt;/p&gt;&lt;/body&gt;&lt;/html&gt;</property>
      <property name="shotcut:animIn">00:00:00.000</property>
      <property name="shotcut:animOut">00:00:00.000</property>
      <property name="disable">0</property>
    </filter>
    <filter id="filter_producer_quote_1_fadeIn" out="00:00:14.520">
      <property name="start">1</property>
      <property name="level">1</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeInBrightness</property>
      <property name="alpha">00:00:00.000=0;00:00:01.000=1</property>
      <property name="shotcut:animIn">00:00:01.000</property>
    </filter>
    <filter id="filter_producer_quote_1_fadeOut" out="00:00:14.520">
      <property name="start">1</property>
      <property name="level">1</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeOutBrightness</property>
      <property name="alpha">00:00:13.520=1;00:00:14.520=0</property>
      <property name="shotcut:animOut">00:00:01.000</property>
    </filter>
  </producer>
  <producer id="producer_quote_2" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">#00000000</property>
    <property name="aspect_ratio">1</property>
    <property name="mlt_service">color</property>
    <property name="mlt_image_format">rgba</property>
    <property name="shotcut:caption">transparent</property>
    <property name="shotcut:detail">transparent</property>
    <property name="ignore_points">0</property>
    <property name="xml">was here</property>
    <property name="seekable">1</property>
    <filter id="producer_filter_producer_quote_2" out="00:00:14.520">
      <property name="argument"/>
      <property name="geometry">902 373 904 546 1</property>
      <property name="family">Sans</property>
      <property name="size">48</property>
      <property name="weight">400</property>
      <property name="style">normal</property>
      <property name="fgcolour">0x000000ff</property>
      <property name="bgcolour">#00000000</property>
      <property name="olcolour">0x00000000</property>
      <property name="pad">0</property>
      <property name="halign">left</property>
      <property name="valign">top</property>
      <property name="outline">0</property>
      <property name="pixel_ratio">1</property>
      <property name="mlt_service">qtext</property>
      <property name="shotcut:filter">richText</property>
      <property name="html">&lt;!DOCTYPE HTML PUBLIC &quot;-//W3C//DTD HTML 4.0//EN&quot; &quot;http://www.w3.org/TR/REC-html40/strict.dtd&quot;&gt;
&lt;html&gt;&lt;head&gt;&lt;meta name=&quot;qrichtext&quot; content=&quot;1&quot; /&gt;&lt;meta charset=&quot;utf-8&quot; /&gt;&lt;style type=&quot;text/css&quot;&gt;
p, li { white-space: pre-wrap; }
hr { height: 1px; border-width: 0; }
li.unchecked::marker { content: &quot;\2610&quot;; }
li.checked::marker { content: &quot;\2612&quot;; }
&lt;/style&gt;&lt;/head&gt;&lt;body style=&quot; font-family:'.AppleSystemUIFont'; font-size:13pt; font-weight:400; font-style:normal;&quot;&gt;
&lt;p align=&quot;center&quot; style=&quot; margin-top:0px; margin-bottom:0px; margin-left:0px; margin-right:0px; -qt-block-indent:0; text-indent:0px; line-height:117%;&quot;&gt;&lt;span style=&quot; font-family:'Mate SC'; font-size:41pt; font-weight:600; color:#ffffff;&quot;&gt;“quote text 3”&lt;/p&gt;&lt;/body&gt;&lt;/html&gt;</property>
      <property name="shotcut:animIn">00:00:00.000</property>
      <property name="shotcut:animOut">00:00:00.000</property>
      <property name="disable">0</property>
    </filter>
    <filter id="filter_producer_quote_2_fadeIn" out="00:00:14.520">
      <property name="start">1</property>
      <property name="level">1</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeInBrightness</property>
      <property name="alpha">00:00:00.000=0;00:00:01.000=1</property>
      <property name="shotcut:animIn">00:00:01.000</property>
    </filter>
    <filter id="filter_producer_quote_2_fadeOut" out="00:00:14.520">
      <property name="start">1</property>
      <property name="level">1</property>
      <property name="mlt_service">brightness</property>
      <property name="shotcut:filter">fadeOutBrightness</property>
      <property name="alpha">00:00:13.520=1;00:00:14.520=0</property>
      <property name="shotcut:animOut">00:00:01.000</property>
    </filter>
  </producer>

  <playlist id="playlist5">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Quote Text</property>
    <blank length="00:00:03.500"/>
    <entry producer="producer_quote_0" in="00:00:00.000" out="00:00:23.440"/>
    <blank length="00:00:02.000"/>
    <entry producer="producer_quote_1" in="00:00:00.000" out="00:00:14.520"/>
    <blank length="00:00:02.000"/>
    <entry producer="producer_quote_2" in="00:00:00.000" out="00:00:14.520"/>
  </playlist>
  <!-- Second track for overlaid background music -->
  <chain id="channel_bgmusic_1" out="00:00:30.000">
    <property name="length">00:00:30.000</property>
    <property name="eof">pause</property>
    <property name="resource">1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3</property>
    <property name="mlt_service">avformat-novalidate</property>
    <property name="seekable">1</property>
    <property name="audio_index">0</property>
    <property name="video_index">-1</property>
    <property name="mute_on_pause">0</property>
    <property name="video_delay">0</property>
    <property name="shotcut:caption">1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3</property>
    <property name="shotcut:producer">avformat-novalidate</property>
    <property name="xml">was here</property>
    <filter id="filter_channel_bgmusic_1_fadeOut" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">00:00:18.000=0;00:00:20.000=-60</property>
      <property name="mlt_service">volume</property>
      <property name="shotcut:filter">fadeOutVolume</property>
      <property name="shotcut:animOut">00:00:02.000</property>
      <property name="disable">0</property>
    </filter>
    <filter id="filter_channel_bgmusic_1_fadeIn" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">00:00:00.000=-60;00:00:02.000=0</property>
      <property name="mlt_service">volume</property>
      <property name="shotcut:filter">fadeInVolume</property>
      <property name="shotcut:animIn">00:00:02.000</property>
    </filter>
    <filter id="filter_channel_bgmusic_1_volume" out="00:00:20.000">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">-11.1</property>
      <property name="mlt_service">volume</property>
    </filter>
  </chain>
  <chain id="channel_bgmusic_3" out="00:00:30.000">
    <property name="length">00:00:30.000</property>
    <property name="eof">pause</property>
    <property name="resource">1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3</property>
    <property name="mlt_service">avformat-novalidate</property>
    <property name="seekable">1</property>
    <property name="audio_index">0</property>
    <property name="video_index">-1</property>
    <property name="mute_on_pause">0</property>
    <property name="video_delay">0</property>
    <property name="shotcut:caption">1-Intro to Inspiration - Calm Inspiring Background Music (Creative Commons) [TubeRipper.com].mp3</property>
    <property name="shotcut:producer">avformat-novalidate</property>
    <property name="xml">was here</property>
    <filter id="filter_channel_bgmusic_3_fadeOut" out="00:00:07.480">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">00:00:05.480=0;00:00:07.480=-60</property>
      <property name="mlt_service">volume</property>
      <property name="shotcut:filter">fadeOutVolume</property>
      <property name="shotcut:animOut">00:00:02.000</property>
      <property name="disable">0</property>
    </filter>
    <filter id="filter_channel_bgmusic_3_fadeIn" out="00:00:07.480">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">00:00:00.000=-60;00:00:02.000=0</property>
      <property name="mlt_service">volume</property>
      <property name="shotcut:filter">fadeInVolume</property>
      <property name="shotcut:animIn">00:00:02.000</property>
    </filter>
    <filter id="filter_channel_bgmusic_3_volume" out="00:00:07.480">
      <property name="window">75</property>
      <property name="max_gain">20dB</property>
      <property name="level">-11.1</property>
      <property name="mlt_service">volume</property>
    </filter>
  </chain>

  <playlist id="playlist_bgmusic_2">
    <property name="shotcut:audio">1</property>
    <property name="shotcut:name">Background Music #2</property>
    <blank length="00:00:18.000"/>
    <entry producer="channel_bgmusic_1" in="00:00:00.000" out="00:00:20.000"/>
    <blank length="00:00:16.000"/>
    <entry producer="channel_bgmusic_3" in="00:00:00.000" out="00:00:07.480"/>
  </playlist>
  <producer id="producer_heading" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">#00000000</property>
    <property name="aspect_ratio">1</property>
    <property name="mlt_service">color</property>
    <property name="mlt_image_format">rgba</property>
    <property name="shotcut:caption">transparent</property>
    <property name="shotcut:detail">transparent</property>
    <property name="ignore_points">0</property>
    <property name="xml">was here</property>
    <property name="seekable">1</property>
    <filter id="filter104" out="00:00:56.120">
      <property name="argument"/>
      <property name="geometry">868 81 962 232 1</property>
      <property name="family">Sans</property>
      <property name="size">48</property>
      <property name="weight">400</property>
      <property name="style">normal</property>
      <property name="fgcolour">0x000000ff</property>
      <property name="bgcolour">#00000000</property>
      <property name="olcolour">0x00000000</property>
      <property name="pad">0</property>
      <property name="halign">left</property>
      <property name="valign">top</property>
      <property name="outline">0</property>
      <property name="pixel_ratio">1</property>
      <property name="mlt_service">qtext</property>
      <property name="shotcut:filter">richText</property>
      <property name="html">&lt;!DOCTYPE HTML PUBLIC &quot;-//W3C//DTD HTML 4.0//EN&quot; &quot;http://www.w3.org/TR/REC-html40/strict.dtd&quot;&gt;
&lt;html&gt;&lt;head&gt;&lt;meta name=&quot;qrichtext&quot; content=&quot;1&quot; /&gt;&lt;meta charset=&quot;utf-8&quot; /&gt;&lt;style type=&quot;text/css&quot;&gt;
p, li { white-space: pre-wrap; }
hr { height: 1px; border-width: 0; }
li.unchecked::marker { content: &quot;\2610&quot;; }
li.checked::marker { content: &quot;\2612&quot;; }
&lt;/style&gt;&lt;/head&gt;&lt;body style=&quot; font-family:'.AppleSystemUIFont'; font-size:13pt; font-weight:400; font-style:normal;&quot;&gt;
&lt;p align=&quot;center&quot; style=&quot; margin-top:0px; margin-bottom:0px; margin-left:0px; margin-right:0px; -qt-block-indent:0; text-indent:0px; line-height:160%;&quot;&gt;&lt;span style=&quot; font-family:'Verdana'; font-size:26pt; font-weight:600; color:#a66d00;&quot;&gt;Title &lt;br&gt; test&lt;/p&gt;&lt;/body&gt;&lt;/html&gt;</property>
      <property name="shotcut:animIn">00:00:00.000</property>
      <property name="shotcut:animOut">00:00:00.000</property>
    </filter>
  </producer>
  <playlist id="playlist_heading">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Heading Text</property>
    <blank length="00:00:03.500"/>
    <entry producer="producer_heading" in="00:00:00.000" out="00:00:56.480"/>
  </playlist>
  <tractor id="tractor0" title="Shotcut version 23.05.14" in="00:00:00.000" out="00:01:07.480">
    <property name="shotcut">1</property>
    <property name="shotcut:projectAudioChannels">2</property>
    <property name="shotcut:projectFolder">0</property>
    <property name="shotcut:scaleFactor">1.96313</property>
    <track producer="playlist_background_color"/>
    <track producer="playlist1"/>
    <track producer="playlist0"/>
    <track producer="playlist_bgmusic_1" hide="video"/>
    <track producer="playlist_outro_image"/>
    <track producer="playlist_overlay"/>
    <track producer="playlist5"/>
    <track producer="playlist_bgmusic_2" hide="video"/>
    <track producer="playlist_heading"/>
    <transition id="transition6">
      <property name="a_track">0</property>
      <property name="b_track">1</property>
      <property name="mlt_service">mix</property>
      <property name="always_active">1</property>
      <property name="sum">1</property>
    </transition>
    <transition id="transition1">
      <property name="a_track">0</property>
      <property name="b_track">2</property>
      <property name="mlt_service">mix</property>
      <property name="always_active">1</property>
      <property name="sum">1</property>
    </transition>
    <transition id="transition0">
      <property name="a_track">1</property>
      <property name="b_track">2</property>
      <property name="version">0.1</property>
      <property name="mlt_service">frei0r.cairoblend</property>
      <property name="1">normal</property>
      <property name="disable">0</property>
    </transition>
    <transition id="transition3">
      <property name="a_track">0</property>
      <property name="b_track">3</property>
      <property name="mlt_service">mix</property>
      <property name="always_active">1</property>
      <property name="sum">1</property>
    </transition>
    <transition id="transition4">
      <property name="a_track">0</property>
      <property name="b_track">4</property>
      <property name="mlt_service">mix</property>
      <property name="always_active">1</property>
      <property name="sum">1</property>
    </transition>
    <transition id="transition2">
      <property name="a_track">1</property>
      <property name="b_track">4</property>
      <property name="version">0.1</property>
      <property name="mlt_service">frei0r.cairoblend</property>
      <property name="1">normal</property>
      <property name="disable">0</property>
    </transition>
    <transition id="transition13">
      <property name="a_track">0</property>
      <property name="b_track">5</property>
      <property name="mlt_service">mix</property>
      <property name="always_active">1</property>
      <property name="sum">1</property>
    </transition>
    <transition id="transition5">
      <property name="a_track">1</property>
      <property name="b_track">5</property>
      <property name="version">0.1</property>
      <property name="mlt_service">frei0r.cairoblend</property>
      <property name="1">normal</property>
      <property name="disable">0</property>
    </transition>
    <transition id="transition7">
      <property name="a_track">0</property>
      <property name="b_track">6</property>
      <property name="mlt_service">mix</property>
      <property name="always_active">1</property>
      <property name="sum">1</property>
    </transition>
    <transition id="transition8">
      <property name="a_track">1</property>
      <property name="b_track">6</property>
      <property name="version">0.1</property>
      <property name="mlt_service">frei0r.cairoblend</property>
      <property name="1">normal</property>
      <property name="disable">0</property>
    </transition>
    <transition id="transition9">
      <property name="a_track">0</property>
      <property name="b_track">7</property>
      <property name="mlt_service">mix</property>
      <property name="always_active">1</property>
      <property name="sum">1</property>
    </transition>
    <transition id="transition10">
      <property name="a_track">0</property>
      <property name="b_track">8</property>
      <property name="mlt_service">mix</property>
      <property name="always_active">1</property>
      <property name="sum">1</property>
    </transition>
    <transition id="transition11">
      <property name="a_track">1</property>
      <property name="b_track">8</property>
      <property name="version">0.1</property>
      <property name="mlt_service">frei0r.cairoblend</property>
      <property name="1">normal</property>
      <property name="disable">0</property>
    </transition>
  </tractor>
</mlt>`

		Expect(err).ToNot(HaveOccurred())
		Expect(xml).To(Equal(expectedXML))
	})
})
