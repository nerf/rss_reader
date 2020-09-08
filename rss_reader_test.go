package rss_reader_test

import (
	"net/http"
	"time"

	. "github.com/nerf/rss_reader"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("RssReader", func() {
	Describe("Parse", func() {
		When("urls set is empty", func() {
			items, err := Parse([]string{})

			It("returns error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("returns empty set", func() {
				Expect(items).Should(BeEmpty())
			})
		})

		When("urls set contains one url", func() {
			var server *ghttp.Server
			var results []RssItem
			var err error

			BeforeEach(func() {
				server = ghttp.NewServer()

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/feed.rss"),
						ghttp.RespondWith(http.StatusOK, `
						<?xml version="1.0" encoding="UTF-8" ?>
						<rss version="2.0">
							<channel>
								<title>Test feed</title>
								<link>http://example.com</link>
								<item>
									<title>Title string</title>
									<description>Description string.</description>
									<link>http://www.example.com/1</link>
									<pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
								</item>
								<item>
									<title>Second Title</title>
									<description>Second Description</description>
									<link>http://www.example.com/2</link>
									<pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
								</item>
							</channel>
						</rss>
						`),
					),
				)

				results, err = Parse([]string{server.URL() + "/feed.rss"})
			})

			AfterEach(func() {
				server.Close()
			})

			It("should return RssItem with data from provided url", func() {
				parsedDate, _ := time.Parse(time.RFC1123Z, "Sun, 06 Sep 2009 16:20:00 +0000")

				Expect(results).Should(HaveLen(2))
				Expect(results[0].Title).To(Equal("Title string"))
				Expect(results[0].Description).To(Equal("Description string."))
				Expect(results[0].Link).To(Equal("http://www.example.com/1"))
				Expect(results[0].PublishDate).To(BeTemporally("==", parsedDate))
				Expect(results[0].Source).To(Equal("Test feed"))
				Expect(results[0].SourceURL).To(Equal(server.URL() + "/feed.rss"))
				Expect(results[1].Title).To(Equal("Second Title"))
			})

			It("error should not be set", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("urls list contains multiple urls", func() {
			var server *ghttp.Server
			var results []RssItem
			var err error

			BeforeEach(func() {
				server = ghttp.NewServer()
				server.RouteToHandler("GET", "/first.rss", ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/first.rss"),
					ghttp.RespondWith(http.StatusOK, `
						<?xml version="1.0" encoding="UTF-8" ?>
						<rss version="2.0">
							<channel>
								<title>First feed</title>
								<link>http://first.com</link>
								<item>
									<title>First Title</title>
									<description>First Description</description>
									<link>http://www.example.com/first</link>
									<pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
								</item>
							</channel>
						</rss>
					`),
				))
				server.RouteToHandler("GET", "/second.rss", ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/second.rss"),
					ghttp.RespondWith(http.StatusOK, `
						<?xml version="1.0" encoding="UTF-8" ?>
						<rss version="2.0">
							<channel>
								<title>Second feed</title>
								<link>http://second.com</link>
								<item>
									<title>Second Title</title>
									<description>Second Description</description>
									<link>http://www.example.com/second</link>
									<pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
								</item>
							</channel>
						</rss>
					`),
				))

				results, err = Parse([]string{server.URL() + "/first.rss", server.URL() + "/second.rss"})
			})

			AfterEach(func() {
				server.Close()
			})

			It("should return RssItem with data from provided urls", func() {
				titles := [2]string{results[0].Title, results[1].Title}

				Expect(results).Should(HaveLen(2))
				Expect(titles).To(ContainElements("First Title", "Second Title"))
			})

			It("error should not be set", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("urls list points to invalid resource", func() {
			var server *ghttp.Server
			var results []RssItem
			var err error

			Context("server response with error", func() {
				BeforeEach(func() {
					server = ghttp.NewServer()

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/feed.rss"),
							ghttp.RespondWith(http.StatusForbidden, ""),
						),
					)

					results, err = Parse([]string{server.URL() + "/feed.rss"})
				})

				AfterEach(func() {
					server.Close()
				})

				It("returns empty response", func() {
					Expect(results).Should(BeEmpty())
				})

				It("not to set error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("server response contains broken payload", func() {
				BeforeEach(func() {
					server = ghttp.NewServer()

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/feed.rss"),
							ghttp.RespondWith(http.StatusOK, "<rss>broken"),
						),
					)

					results, err = Parse([]string{server.URL() + "/feed.rss"})
				})

				AfterEach(func() {
					server.Close()
				})

				It("returns empty response", func() {
					Expect(results).Should(BeEmpty())
				})

				It("not to set error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})
})
