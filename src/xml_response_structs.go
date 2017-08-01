package main

type sitemap struct {
	Lat  string `xml:"Latitude"`
	Lng  string `xml:"Longitude"`
	Name string `xml:"Name,attr"`
}

type sitemapList struct {
	Sites []sitemap `xml:"Site"`
}

type xmlErrorResponse struct {
	Error string `xml:"Error"`
}

type rainSite struct {
	Measurement struct {
		Data struct {
			Intervals []rainFallInterval `xml:"E"`
		} `xml:"Data"`
	} `xml:"Measurement"`
}

type rainFallInterval struct {
	Time     string `xml:"T"`
	RainFall string `xml:"I1"`
}
