package templates

var TemplateQuotes = `<?xml version="1.0" standalone="no"?>
<mlt LC_NUMERIC="C" version="7.17.0" title="Shotcut version 23.05.14" producer="main_bin">
  <profile description="PAL 4:3 DV or DVD" width="1920" height="1080" progressive="1" sample_aspect_num="1" sample_aspect_den="1" display_aspect_num="16" display_aspect_den="9" frame_rate_num="25" frame_rate_den="1" colorspace="709"/>
  <playlist id="main_bin">
    <property name="xml_retain">1</property>
  </playlist>
  <producer id="producer_background_color" in="00:00:00.000" out="{{.TotalVideoDuration}}">
    <property name="length">{{.TotalVideoDuration}}</property>
    <property name="eof">pause</property>
    <property name="resource">0</property>
    <property name="aspect_ratio">1</property>
    <property name="mlt_service">color</property>
    <property name="mlt_image_format">rgba</property>
    <property name="set.test_audio">0</property>
  </producer>
  <playlist id="playlist_background_color">
    <property name="shotcut:name">Background Color</property>
    <entry producer="producer_background_color" in="00:00:00.000" out="{{.TotalVideoDuration}}"/>
  </playlist>
  <producer id="bg_image_producer" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">{{.BackgroundImagePath}}</property>
    <property name="ttl">1</property>
    <property name="aspect_ratio">1</property>
    <property name="progressive">1</property>
    <property name="seekable">1</property>
    <property name="format">1</property>
    <property name="mlt_service">qimage</property>
    <property name="creation_time">2023-05-23T13:36:16</property>
    <property name="shotcut:caption">{{.BackgroundImageName}}</property>
    <property name="xml">was here</property>
  </producer>
  <playlist id="playlist0">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Background Image</property>
    <entry producer="bg_image_producer" in="00:00:00.000" out="{{.TotalVideoDuration}}"/>
  </playlist>
{{.QuoteAudioChannelsXML}}
  <playlist id="playlist1">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Quote Audio</property>
    <blank length="{{.IntroDuration}}"/>
{{.QuoteAudioPlaylistEntriesXML}}
  </playlist>
  <!-- First track for background music -->
{{.BgMusicChainsFirstTrackXML}}
  <playlist id="playlist_bgmusic_1">
    <property name="shotcut:audio">1</property>
    <property name="shotcut:name">Background Music #1</property>
{{.BgMusicFirstTrackPlaylistEntriesXML}}
  </playlist>
  <producer id="producer_intro_image" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">{{.IntroImagePath}}</property>
    <property name="ttl">1</property>
    <property name="aspect_ratio">1</property>
    <property name="progressive">1</property>
    <property name="seekable">1</property>
    <property name="format">1</property>
    <property name="mlt_service">qimage</property>
    <property name="creation_time">2023-05-23T00:26:22</property>
    <property name="xml">was here</property>
    <property name="shotcut:caption">{{.IntroImageName}}</property>
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
    <property name="resource">{{.OutroImagePath}}</property>
    <property name="ttl">1</property>
    <property name="aspect_ratio">1</property>
    <property name="progressive">1</property>
    <property name="seekable">1</property>
    <property name="format">1</property>
    <property name="mlt_service">qimage</property>
    <property name="creation_time">2023-05-23T00:42:12</property>
    <property name="shotcut:caption">{{.OutroImageName}}</property>
    <property name="xml">was here</property>
  </producer>
  <playlist id="playlist_outro_image">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Intro and Outro</property>
    <property name="shotcut:lock">0</property>
    <entry producer="producer_intro_image" in="00:00:00.000" out="{{.IntroDuration}}"/>
    <blank length="{{.IntroOutroGapDuration}}"/>
    <entry producer="producer_outro_image" in="00:00:00.000" out="{{.OutroDuration}}"/>
  </playlist>
  <producer id="producer_post_outro_image" in="00:00:00.000" out="03:59:59.960">
    <property name="length">04:00:00.000</property>
    <property name="eof">pause</property>
    <property name="resource">{{.OutroOverlayImagePath}}</property>
    <property name="ttl">1</property>
    <property name="aspect_ratio">1</property>
    <property name="progressive">1</property>
    <property name="seekable">1</property>
    <property name="format">1</property>
    <property name="mlt_service">qimage</property>
    <property name="creation_time">2023-05-23T00:41:13</property>
    <property name="shotcut:caption">{{.OutroOverlayImageName}}</property>
    <property name="xml">was here</property>
    <filter id="filter8" out="{{.OutroOverlayDuration}}">
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
    <blank length="{{.OutroOverlayGap}}"/>
    <entry producer="producer_post_outro_image" in="00:00:00.000" out="{{.OutroOverlayDuration}}"/>
  </playlist>
{{.QuoteTextProducersXML}}
  <playlist id="playlist5">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Quote Text</property>
    <blank length="{{.IntroDuration}}"/>
{{.QuoteTextProducerEntriesXML}}
  </playlist>
  <!-- Second track for overlaid background music -->
{{.BgMusicChainsSecondTrackXML}}
  <playlist id="playlist_bgmusic_2">
    <property name="shotcut:audio">1</property>
    <property name="shotcut:name">Background Music #2</property>
{{.BgMusicSecondTrackPlaylistEntriesXML}}
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
&lt;p align=&quot;center&quot; style=&quot; margin-top:0px; margin-bottom:0px; margin-left:0px; margin-right:0px; -qt-block-indent:0; text-indent:0px; line-height:160%;&quot;&gt;&lt;span style=&quot; font-family:'Verdana'; font-size:26pt; font-weight:600; color:#a66d00;&quot;&gt;{{.HeadingHtml}}&lt;/p&gt;&lt;/body&gt;&lt;/html&gt;</property>
      <property name="shotcut:animIn">00:00:00.000</property>
      <property name="shotcut:animOut">00:00:00.000</property>
    </filter>
  </producer>
  <playlist id="playlist_heading">
    <property name="shotcut:video">1</property>
    <property name="shotcut:name">Heading Text</property>
    <blank length="00:00:03.500"/>
    <entry producer="producer_heading" in="00:00:00.000" out="{{.HeadingDuration}}"/>
  </playlist>
  <tractor id="tractor0" title="Shotcut version 23.05.14" in="00:00:00.000" out="{{.TotalVideoDuration}}">
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

type QuotesTemplateParams struct {
	BackgroundImagePath string
	BackgroundImageName string
	BackgroundImageHash string

	IntroImagePath string
	IntroImageName string
	IntroImageHash string
	IntroDuration  string

	OutroImagePath        string
	OutroImageName        string
	OutroImageHash        string
	OutroDuration         string
	IntroOutroGapDuration string

	OutroOverlayImagePath string
	OutroOverlayImageName string
	OutroOverlayImageHash string
	OutroOverlayGap       string
	OutroOverlayDuration  string

	TotalVideoDuration string

	HeadingDuration string
	HeadingHtml     string

	QuoteAudioChannelsXML        string
	QuoteAudioPlaylistEntriesXML string

	QuoteTextProducersXML       string
	QuoteTextProducerEntriesXML string

	BgMusicChainsFirstTrackXML           string
	BgMusicChainsSecondTrackXML          string
	BgMusicFirstTrackPlaylistEntriesXML  string
	BgMusicSecondTrackPlaylistEntriesXML string
}

func QuotesTemplateParamsNew(
	backgroundImagePath string,
	backgroundImageName string,

	introImagePath string,
	introImageName string,
	introDuration string,

	outroImagePath string,
	outroImageName string,
	outroDuration string,
	introOutroGapDuration string,

	outroOverlayImagePath string,
	outroOverlayImageName string,
	outroOverlayGap string,
	outroOverlayDuration string,

	totalVideoDuration string,

	headingDuration string,
	headingHtml string,

	quoteAudioChannelsXML string,
	quoteAudioPlaylistEntriesXML string,

	quoteTextProducersXML string,
	quoteTextProducerEntriesXML string,

	bgMusicChainsFirstTrackXML string,
	bgMusicChainsSecondTrackXML string,
	bgMusicFirstTrackPlaylistEntriesXML string,
	bgMusicSecondTrackPlaylistEntriesXML string,
) QuotesTemplateParams {
	return QuotesTemplateParams{
		BackgroundImagePath: backgroundImagePath,
		BackgroundImageName: backgroundImageName,

		IntroImagePath: introImagePath,
		IntroImageName: introImageName,
		IntroDuration:  introDuration,

		OutroImagePath:        outroImagePath,
		OutroImageName:        outroImageName,
		OutroDuration:         outroDuration,
		IntroOutroGapDuration: introOutroGapDuration,

		OutroOverlayImagePath: outroOverlayImagePath,
		OutroOverlayImageName: outroOverlayImageName,
		OutroOverlayGap:       outroOverlayGap,
		OutroOverlayDuration:  outroOverlayDuration,

		TotalVideoDuration: totalVideoDuration,

		HeadingDuration: headingDuration,
		HeadingHtml:     headingHtml,

		QuoteAudioChannelsXML:        quoteAudioChannelsXML,
		QuoteAudioPlaylistEntriesXML: quoteAudioPlaylistEntriesXML,

		QuoteTextProducersXML:       quoteTextProducersXML,
		QuoteTextProducerEntriesXML: quoteTextProducerEntriesXML,

		BgMusicChainsFirstTrackXML:           bgMusicChainsFirstTrackXML,
		BgMusicChainsSecondTrackXML:          bgMusicChainsSecondTrackXML,
		BgMusicFirstTrackPlaylistEntriesXML:  bgMusicFirstTrackPlaylistEntriesXML,
		BgMusicSecondTrackPlaylistEntriesXML: bgMusicSecondTrackPlaylistEntriesXML,
	}
}
