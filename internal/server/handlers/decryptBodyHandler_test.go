package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"

	"github.com/gin-gonic/gin"
)

// CheckMetricExistenceHandler checks metrics has value set
func TestNewDecryptBodyHandler(t *testing.T) {
	privK, pubK := rsacrypt.GenerateKeyPair(2048)
	payload := `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur tempor non turpis a pretium. Vivamus dapibus pellentesque odio, in sagittis tellus tincidunt eget. Morbi ornare, elit at vestibulum vestibulum, quam ex pharetra velit, dictum maximus mi lacus in risus. Integer velit ligula, suscipit ac cursus at, blandit eget neque. Duis quis tristique justo, eget pellentesque leo. Nunc nisi orci, dignissim at libero sed, fringilla tempor magna. Vivamus nibh est, condimentum ut viverra at, condimentum eget urna. Praesent sed neque pulvinar, vestibulum libero id, euismod mauris. Pellentesque pharetra lectus sed est cursus, nec mattis dolor ornare. Phasellus bibendum bibendum interdum. Fusce efficitur ultricies dignissim. Sed eu eros fermentum, vestibulum orci in, tincidunt orci.

In hac habitasse platea dictumst. In maximus suscipit lacinia. Aliquam rhoncus non nisl nec porttitor. Aenean interdum lacus nec eros dignissim, a fermentum felis imperdiet. Morbi a vulputate arcu. Cras porta, risus et auctor rhoncus, ligula libero sollicitudin enim, a dictum orci quam at neque. Fusce sollicitudin magna a orci feugiat, ac scelerisque libero volutpat. Donec egestas velit id lectus luctus porta. Phasellus dignissim tortor porta dolor accumsan, sit amet semper dolor tristique. Nulla purus ex, viverra in lorem a, malesuada varius enim. Pellentesque luctus vulputate urna et mattis. Sed in leo vel dui vehicula semper a a urna.

Praesent venenatis in ex sit amet pellentesque. Ut non ex vitae nisi venenatis iaculis sit amet ut augue. Ut vel gravida leo, nec pulvinar nisi. Aenean ornare ex eget finibus laoreet. In volutpat sit amet quam vel convallis. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec convallis accumsan diam ac aliquet. Sed lobortis tempor erat sit amet egestas. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Sed at tincidunt felis. Sed nisl leo, mollis et iaculis vel, fringilla vitae tellus. Donec et finibus urna. Maecenas congue augue a euismod porttitor. Nullam gravida, erat sed fringilla varius, risus nunc facilisis elit, in luctus nisi sem a ante. In hac habitasse platea dictumst. Nam a accumsan odio, tempus pellentesque leo.

Vestibulum volutpat semper nulla, nec efficitur nulla iaculis ut. Etiam ornare odio quis ipsum dapibus, tincidunt rutrum est mollis. Suspendisse iaculis laoreet turpis malesuada pharetra. Praesent luctus viverra tincidunt. Aenean vel tellus hendrerit, viverra dui in, scelerisque nibh. Nullam lobortis est sed tempus fringilla. Donec dignissim metus mi. Ut facilisis finibus lacus, eget consequat augue rutrum sed. Maecenas vitae mauris ante. Cras feugiat ligula a tellus ultricies, vitae hendrerit risus pharetra. Proin id quam ac augue pretium tincidunt in a tellus. Maecenas efficitur diam augue, eget pellentesque eros porta eget.

Sed laoreet lobortis arcu quis congue. Curabitur dapibus hendrerit quam, sit amet sollicitudin libero euismod id. Suspendisse vel commodo nibh. Vivamus ullamcorper, elit sit amet accumsan consectetur, sapien metus aliquam nisl, sed imperdiet leo mauris id mi. Donec eget sapien fermentum, blandit metus eget, lobortis sapien. Proin volutpat, orci sed vulputate aliquet, sapien nisl efficitur enim, non consequat ex diam semper tortor. Vivamus est felis, imperdiet id metus sit amet, aliquam luctus augue. Sed commodo odio ligula, quis pharetra leo eleifend nec. Integer ut tincidunt metus. Aenean vestibulum hendrerit lectus non faucibus. Nulla consectetur mi in neque mollis, quis mollis ante tempus. Nam vitae libero accumsan, sodales neque vitae, tristique velit. Aliquam at odio dolor. Donec nisl enim, luctus vitae dui quis, imperdiet gravida mi. 
`
	crypted, err := rsacrypt.EncryptWithPublicKey([]byte(payload), pubK)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}
	tests := []struct {
		name        string
		requestBody []byte
		withError   bool
		wantErr     bool
	}{
		{
			name:        "fail",
			requestBody: []byte(payload),
			wantErr:     true,
		},
		{
			name:        "success",
			requestBody: crypted,
		},
	}
	h := NewDecryptBodyHandler(privK)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			rc := bytes.NewBuffer(tt.requestBody)
			r := httptest.NewRequest(http.MethodPost, "/update", rc)
			c.Request = r
			h(c)

			if (w.Code != http.StatusOK || len(c.Errors) > 0) != tt.wantErr {
				t.Errorf("NewDecryptBodyHandler() code = %d, errors = %v, wantErr %v", w.Code, len(c.Errors), tt.wantErr)
			}

			b := make([]byte, len(tt.requestBody))
			c.Request.Body.Read(b)
		})
	}
}
