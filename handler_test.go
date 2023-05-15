package main

import (
	"bytes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = BeforeSuite(func() {
	go func() {
		defer GinkgoRecover()
	}()
})

var _ = Describe("Handlers", func() {
	Context("with fresh store", func() {
		var svc Service
		svc = NewStore()
		svc = NewAutoFields(svc)

		h := NewHandler(svc)

		It("should fail on list unknown kind", func() {
			req := httptest.NewRequest(http.MethodGet, "/foo", nil)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusNotFound))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("after creating an item", Ordered, func() {
		var svc Service
		svc = NewStore()
		svc = NewAutoFields(svc)

		h := NewHandler(svc)

		BeforeAll(func() {
			reqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz"}`)

			req := httptest.NewRequest(http.MethodPost, "/foo", reqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusCreated))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("id", "foo1"))
			Expect(rsp).Should(HaveKeyWithValue("kind", "foo"))
			Expect(rsp).Should(HaveKeyWithValue("bar", "baz"))
			Expect(rsp).Should(HaveKey("uuid"))
			Expect(rsp).Should(HaveKey("createdAt"))
			Expect(rsp).Should(HaveKey("updatedAt"))
		})

		It("should have 1 item", func() {
			req := httptest.NewRequest(http.MethodGet, "/foo", nil)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var rsp ListResponse

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp.Items).Should(HaveLen(1))
		})

		It("should be able to read", func() {
			req := httptest.NewRequest(http.MethodGet, "/foo/foo1", nil)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("id", "foo1"))
			Expect(rsp).Should(HaveKeyWithValue("kind", "foo"))
			Expect(rsp).Should(HaveKeyWithValue("bar", "baz"))
			Expect(rsp).Should(HaveKey("uuid"))
			Expect(rsp).Should(HaveKey("createdAt"))
			Expect(rsp).Should(HaveKey("updatedAt"))
			Expect(rsp["updatedAt"]).Should(Equal(rsp["createdAt"]))
		})

		It("should fail on read from unknown kind", func() {
			req := httptest.NewRequest(http.MethodGet, "/bar/bar1", nil)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusNotFound))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should not be able to create again", func() {
			reqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz"}`)

			req := httptest.NewRequest(http.MethodPost, "/foo", reqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusConflict))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("message", "Item exists"))
		})
	})

	Context("on creating an invalid item", Ordered, func() {
		var svc Service
		svc = NewStore()
		svc = NewAutoFields(svc)

		h := NewHandler(svc)

		reqBody := bytes.NewBufferString(`{"id":"foo1","bar":INVALID JSON}`)

		req := httptest.NewRequest(http.MethodPost, "/foo", reqBody)

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		res := w.Result()

		defer func() { _ = res.Body.Close() }()

		It("should fail", func() {
			Expect(res.StatusCode).Should(Equal(http.StatusBadRequest))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("message", "Invalid request"))
		})
	})

	Context("after replacing an item", Ordered, func() {
		var svc Service
		svc = NewStore()
		svc = NewAutoFields(svc)

		h := NewHandler(svc)

		BeforeAll(func() {
			createReqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz"}`)

			createReq := httptest.NewRequest(http.MethodPost, "/foo", createReqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, createReq)

			createRes := w.Result()

			defer func() { _ = createRes.Body.Close() }()

			Expect(createRes.StatusCode).Should(Equal(http.StatusCreated))

			updateReqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz2", "bool": true}`)

			updateReq := httptest.NewRequest(http.MethodPut, "/foo/foo1", updateReqBody)

			w2 := httptest.NewRecorder()

			h.ServeHTTP(w2, updateReq)

			updateRes := w2.Result()

			defer func() { _ = updateRes.Body.Close() }()

			Expect(updateRes.StatusCode).Should(Equal(http.StatusNoContent))
		})

		It("should be updated", func() {
			req := httptest.NewRequest(http.MethodGet, "/foo/foo1", nil)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("id", "foo1"))
			Expect(rsp).Should(HaveKeyWithValue("kind", "foo"))
			Expect(rsp).Should(HaveKeyWithValue("bar", "baz2"))
			Expect(rsp).Should(HaveKeyWithValue("bool", true))
			Expect(rsp["updatedAt"]).ShouldNot(Equal(rsp["createdAt"]))
		})
	})

	Context("on replacing an item with invalid data", Ordered, func() {
		var svc Service
		svc = NewStore()
		svc = NewAutoFields(svc)

		h := NewHandler(svc)

		BeforeAll(func() {
			createReqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz"}`)

			createReq := httptest.NewRequest(http.MethodPost, "/foo", createReqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, createReq)

			createRes := w.Result()

			defer func() { _ = createRes.Body.Close() }()

			Expect(createRes.StatusCode).Should(Equal(http.StatusCreated))
		})

		It("should fail", func() {
			reqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz2", BAD JSON}`)

			req := httptest.NewRequest(http.MethodPut, "/foo/foo1", reqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusBadRequest))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("message", "Invalid request"))
		})
	})

	Context("on replacing an not existed item", Ordered, func() {
		var svc Service
		svc = NewStore()
		svc = NewAutoFields(svc)

		h := NewHandler(svc)

		BeforeAll(func() {
			createReqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz"}`)

			createReq := httptest.NewRequest(http.MethodPost, "/foo", createReqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, createReq)

			createRes := w.Result()

			defer func() { _ = createRes.Body.Close() }()

			Expect(createRes.StatusCode).Should(Equal(http.StatusCreated))
		})

		It("should fail", func() {
			reqBody := bytes.NewBufferString(`{"id":"foo2","bar":"baz2", "bool": true}`)

			req := httptest.NewRequest(http.MethodPut, "/foo/foo2", reqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusNotFound))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("message", "Item not found"))
		})

		It("should fail on replace in unknown kind", func() {
			reqBody := bytes.NewBufferString(`{"id":"foo2","bar":"baz2", "bool": true}`)

			req := httptest.NewRequest(http.MethodPut, "/bar/bar2", reqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusNotFound))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("message", "Invalid kind"))
		})
	})

	Context("after deleting an item", Ordered, func() {
		var svc Service
		svc = NewStore()
		svc = NewAutoFields(svc)

		h := NewHandler(svc)

		BeforeAll(func() {
			createReqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz"}`)

			createReq := httptest.NewRequest(http.MethodPost, "/foo", createReqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, createReq)

			createRes := w.Result()

			defer func() { _ = createRes.Body.Close() }()

			Expect(createRes.StatusCode).Should(Equal(http.StatusCreated))

			deleteReq := httptest.NewRequest(http.MethodDelete, "/foo/foo1", nil)

			w2 := httptest.NewRecorder()

			h.ServeHTTP(w2, deleteReq)

			deleteRes := w2.Result()

			defer func() { _ = deleteRes.Body.Close() }()

			Expect(deleteRes.StatusCode).Should(Equal(http.StatusNoContent))
		})

		It("should be not found", func() {
			req := httptest.NewRequest(http.MethodGet, "/foo/foo1", nil)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusNotFound))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("message", "Item not found"))
		})
	})

	Context("on deleting an not existed item", Ordered, func() {
		var svc Service
		svc = NewStore()
		svc = NewAutoFields(svc)

		h := NewHandler(svc)

		BeforeAll(func() {
			createReqBody := bytes.NewBufferString(`{"id":"foo1","bar":"baz"}`)

			createReq := httptest.NewRequest(http.MethodPost, "/foo", createReqBody)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, createReq)

			createRes := w.Result()

			defer func() { _ = createRes.Body.Close() }()

			Expect(createRes.StatusCode).Should(Equal(http.StatusCreated))
		})

		It("should fail", func() {
			req := httptest.NewRequest(http.MethodDelete, "/foo/foo2", nil)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusNotFound))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rsp).Should(HaveKeyWithValue("message", "Item not found"))
		})

		It("should fail on delete from unknown kind", func() {
			req := httptest.NewRequest(http.MethodDelete, "/bar/bar1", nil)

			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()

			defer func() { _ = res.Body.Close() }()

			Expect(res.StatusCode).Should(Equal(http.StatusNotFound))

			var rsp map[string]interface{}

			err := json.NewDecoder(res.Body).Decode(&rsp)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
