package hetznerapi

import (
	"fmt"
	"log"
	// "encoding/json"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"context"
	"golang.org/x/crypto/ssh"
	"github.com/google/uuid"
	"strings"
)

type account struct {
    apiKey			string
    hetznerAPI		*hcloud.Client
}

func New(apiKey string) account {  
    client := hcloud.NewClient(hcloud.WithToken(apiKey))

    _, err := client.Server.All(context.Background())
    // fmt.Println(stt)

	if err != nil {
		log.Fatal("ERR:",err)
	}
	a := account {apiKey: apiKey, hetznerAPI: client}

    return a
}

func (a account) ListAllServers() []*hcloud.Server { 
    servers, err := a.hetznerAPI.Server.All(context.Background())

	if err != nil {
		log.Fatal("ERR:",err)
	}

    return servers
}

func (a account) GetSSHKeyFingerprint (pem string) string {
	pubKeyBytes := []byte(pem)
	pk, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
    if err != nil {
        log.Fatal(err)
        return ""
    }

    // Get the fingerprint
    f := ssh.FingerprintLegacyMD5(pk)

    return f
}

func (a account) GetSSHKeyIdByFingerprint(fingerprint string) int { 
    sshKey, _, err := a.hetznerAPI.SSHKey.GetByFingerprint(context.Background(), fingerprint)
	if err != nil {
		log.Fatal("SSHKey.GetByFingerprint failed: %s", err)
		return -1
	}
	if sshKey == nil {
		log.Fatal("no SSH key")
		return -1
	}

    return sshKey.ID
}

func (a account) CreateSSHKey(name string, pem string) int { 
    opts := hcloud.SSHKeyCreateOpts{
		Name:      name,
		PublicKey: pem,
	}
	sshKey, _, err := a.hetznerAPI.SSHKey.Create(context.Background(), opts)
	if err != nil {
		log.Fatal("SSHKey.Get failed: %s", err)
		return -1
	}
	
    return sshKey.ID
}

func (a account) SSHKeyIdGetOrCreate(pem string) int {
	f := a.GetSSHKeyFingerprint(pem)
	key := a.GetSSHKeyIdByFingerprint(f)
	if key != -1 {
		// fmt.Println("Existed")
		return key
	}
	// fmt.Println("Creating...")
	id := a.CreateSSHKey(f, pem)
	if id != -1 {
		// fmt.Println("Created")
		return id
	}
	// fmt.Println("Failed")
	return -1
}

func (a account) GetAllLocation() { 
	objs, err := a.hetznerAPI.Location.All(context.Background())
	if err != nil {
		log.Fatalf("X.List failed: %s", err)
	}
	for _, element := range objs {
		fmt.Println(element.ID, element.Name)
	}
	// return -1
}
func (a account) GetAllImages() { 
	objs, err := a.hetznerAPI.Image.All(context.Background())
	if err != nil {
		log.Fatalf("X.List failed: %s", err)
	}
	for _, element := range objs {
		fmt.Println(element.ID, element.Name)
	}
	// return -1
}

func (a account) GetServerTypeIdByName(name string) int { 
	serverTypes, err := a.hetznerAPI.ServerType.All(context.Background())
	if err != nil {
		log.Fatalf("ServerTypes.List failed: %s", err)
	}
	for _, element := range serverTypes {
		if element.Name == name {
			return element.ID
		}
	}
	return -1
}
func (a account) CreateServer(sshKeyId int, svType string, svImage string, svLocation string) int { 
    result, _, err := a.hetznerAPI.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:       strings.Replace(uuid.New().String(),"-",".",-1)+".local",
		ServerType: &hcloud.ServerType{Name: svType},
		Image:      &hcloud.Image{Name: svImage},
		Location:   &hcloud.Location{Name: svLocation},
		UserData:   "#cloud-config\nruncmd:\n- [yum, update, -y]\n",
		SSHKeys: []*hcloud.SSHKey{
			{ID: sshKeyId},
		},
	})
	if err != nil {
		log.Fatalf("Server.Create failed: %s", err)
	}
	if result.Server == nil {
		log.Fatal("no server")
	}
	// if result.RootPassword != "" {
	// 	log.Fatalf("expected no root password, got: %v", result.RootPassword)
	// }
	return result.Server.ID
}

func (a account) DeleteServer(serverId int) bool { 
    _, err := a.hetznerAPI.Server.Delete(context.Background(), &hcloud.Server{ID: serverId})
	if err != nil {
		log.Fatalf("Server.Delete failed: %s", err)
		return false
	}
	return true
}
func (a account) GetServerById(id int) hcloud.Server { 
    server, _, err := a.hetznerAPI.Server.GetByID(context.Background(), id)

	if err != nil {
		log.Fatal("ERR:",err)
	}

    return *server
}