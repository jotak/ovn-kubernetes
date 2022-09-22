/*
 * HCS API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 2.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package hcsschema

type Devices struct {
	ComPorts map[string]ComPort `json:"ComPorts,omitempty"`

	Scsi map[string]Scsi `json:"Scsi,omitempty"`

	VirtualPMem *VirtualPMemController `json:"VirtualPMem,omitempty"`

	NetworkAdapters map[string]NetworkAdapter `json:"NetworkAdapters,omitempty"`

	VideoMonitor *VideoMonitor `json:"VideoMonitor,omitempty"`

	Keyboard *Keyboard `json:"Keyboard,omitempty"`

	Mouse *Mouse `json:"Mouse,omitempty"`

	HvSocket *HvSocket2 `json:"HvSocket,omitempty"`

	EnhancedModeVideo *EnhancedModeVideo `json:"EnhancedModeVideo,omitempty"`

	GuestCrashReporting *GuestCrashReporting `json:"GuestCrashReporting,omitempty"`

	VirtualSmb *VirtualSmb `json:"VirtualSmb,omitempty"`

	Plan9 *Plan9 `json:"Plan9,omitempty"`

	Battery *Battery `json:"Battery,omitempty"`

	FlexibleIov map[string]FlexibleIoDevice `json:"FlexibleIov,omitempty"`

	SharedMemory *SharedMemoryConfiguration `json:"SharedMemory,omitempty"`

	// TODO: This is pre-release support in schema 2.3. Need to add build number
	// docs when a public build with this is out.
	VirtualPci map[string]VirtualPciDevice `json:",omitempty"`
}