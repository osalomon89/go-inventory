package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/osalomon89/go-inventory/internal/domain"
	"github.com/osalomon89/go-inventory/internal/server"
	"github.com/osalomon89/go-inventory/internal/test/mocks"
)

var _ = Describe("Server", func() {
	var r *http.Request
	var w *httptest.ResponseRecorder
	var bookRepositoryMock *mocks.MockBookRepository
	var mockCtrl *gomock.Controller

	BeforeEach(func() {
		r = httptest.NewRequest(http.MethodGet, "/books", nil)
		w = httptest.NewRecorder()

		mockCtrl = gomock.NewController(GinkgoT())
		bookRepositoryMock = mocks.NewMockBookRepository(mockCtrl)
	})

	AfterEach(func() {
		w.Result().Body.Close()
		mockCtrl.Finish()
	})

	Context("When get request is sent with an valid ID", func() {
		BeforeEach(func() {
			bookRepositoryMock.EXPECT().
				GetBookByID(uint(1)).
				Return(&domain.Book{
					ID:     1,
					Author: "the-author",
					Title:  "the title",
					Price:  5000,
					Isbn:   "abcd1234",
					Stock:  20,
				}, nil).Times(1)

			r = mux.SetURLVars(r, map[string]string{"id": "1"})
		})
		It("Returns a book", func() {
			handler := server.NewHandler(bookRepositoryMock)

			handler.GetBookByID(w, r)

			Expect(w.Result().StatusCode).Should(Equal(http.StatusOK))
			body, err := io.ReadAll(w.Body)
			Expect(err).ShouldNot(HaveOccurred())

			var result server.ResponseInfo
			err = json.Unmarshal(body, &result)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(result.Data).To(Equal(domain.Book{
				ID:     1,
				Author: "the-author",
				Title:  "the title",
				Price:  5000,
				Isbn:   "abcd1234",
				Stock:  20,
			}))
		})
	})
})
