package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"time"

	_ "github.com/lib/pq"
)

const (
	sitemapListURL  = "http://data.hbrc.govt.nz/Envirodata/EMAR.hts?service=Hilltop&Request=SiteList&Location=LatLong"
	rainMapCheckURL = "http://data.hbrc.govt.nz/Envirodata/EMAR.hts?service=Hilltop&request=GetData&Site=%s&Measurement=Rainfall"
	rainAt15Minutes = "http://data.hbrc.govt.nz/Envirodata/EMAR.hts?service=Hilltop&request=GetData&Site=%s&Measurement=Rainfall&From=%s&To=%s&Interval=15%20Minutes&Interpolation=Total"
	months          = 12
)

//The rainMapCheckUrl will return an error if a rain station doesn't
//have any data for time specified

func main() {
	siteMaps := getSitemaps()
	for _, site := range siteMaps.Sites {
		addSitemapToDatabase(site)
	}
}

func addSitemapToDatabase(site sitemap) {
	fmt.Println("/n-/n-Checking Site:", site.Name)
	if rainResponse := checkSiteMapTypeRain(site.Name); rainResponse != "" {
		fmt.Println("--Adding Site:", site.Name)
		id, err := dbAddSite(site)
		if err == nil {
			end := time.Now()
			for month := 0; month < months; month++ {
				start = end.time - month
				rainRecords := getRainRecords(site.Name, start, end)
				for rainfall, _ := range rainRecords {
					dbAddRainRecord(rainfall, id)
				}
				end - month
			}
		}
	} else {
		fmt.Println("--Site", site.Name, "rain response:", rainResponse)
	}
}

func getRainRecords(name string, start, end time.Time) []rainFallInterval {
	var rainSitemap rainSite

	url := fmt.Sprintf(rainAt15Minutes, strings.Replace(name, " ", "%20", -1),
		start.Format("02/01/2006"), end.Format("02/01/2006"))
	xmlByteResponse, err := getContent(url)
	checkErr(err)

	err = xml.Unmarshal(xmlByteResponse, &rainSitemap)
	checkErr(err)

	return rainSitemap.Measurement.Data.Intervals
}

func checkSiteMapTypeRain(name string) string {
	var xmlResponse xmlErrorResponse

	url := fmt.Sprintf(rainMapCheckURL, strings.Replace(name, " ", "%20", -1))
	xmlByteResponse, err := getContent(url)
	checkErr(err)

	err = xml.Unmarshal(xmlByteResponse, &xmlResponse)
	checkErr(err)

	return xmlResponse.Error
}

func getSitemaps() sitemapList {
	var sitemaps sitemapList
	sitemapXMLList, err := getContent(sitemapListURL)
	checkErr(err)
	err = xml.Unmarshal(sitemapXMLList, &sitemaps)
	checkErr(err)
	return sitemaps
}

func getContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
