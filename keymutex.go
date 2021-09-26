import 	"k8s.io/utils/keymutex"



func TestKeyMutex(t *testing.T) {
	km := keymutex.NewHashed(0)

	km.LockKey(cast.ToString("lemmaId"))
	defer func(km keymutex.KeyMutex, id string) {
		_ = km.UnlockKey(id)
	}(km, cast.ToString("lemmaId"))

	// operate on this lemma
}
