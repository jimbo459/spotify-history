package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"

	handlers "github.com/jimbo459/spotify-history/src/internal/handlers"
)

var _ = Describe("Handlers", func() {

	Describe("CallBack Handler", func() {
		When("Successful query is made", func() {
			It("Returns status OK", func() {
				req, err := http.NewRequest("GET", "/callback", nil)
				Expect(err).NotTo(HaveOccurred())

				ResponseRecorder := httptest.NewRecorder()

				handler := http.HandlerFunc(handlers.CallBackHandler)

				handler.ServeHTTP(ResponseRecorder, req)

				Expect(ResponseRecorder.Code).To(Equal(http.StatusOK))

			})
		})
	})

})
