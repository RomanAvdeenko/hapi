package mysql

import (
	"fmt"
	"hapi/internal/app/model"
	"hapi/internal/app/store"

	"github.com/RomanAvdeenko/utils/net"
)

const (
	querySelectUserPackageByUserID = `
	SELECT 
		userpackage_id,
		userpackage_unique_id,
		userpackage_parent,
		userpackage_user_id,
		userpackage_package_id,
		userpackage_start,
		userpackage_stop,
		userpackage_enabled,
		userpackage_lastmodify,
		COALESCE(userpackage_ip, '0'),
		COALESCE(userpackage_ip_network, '0'),
		userpackage_mac,
		userpackage_desc
	FROM userpackage
	WHERE
    	userpackage_start<=Now() AND
    	userpackage_stop>=Now() AND
		userpackage_user_id = ?`
)

type UserPackageRepository struct {
	store *Store
}

func (r *UserPackageRepository) FindByUserID(userid uint) ([]model.UserPackage, error) {
	got := make([]model.UserPackage, 0, 4)

	//fmt.Println(querySelectUserPackageByUserID, userid)

	rows, err := r.store.db.Query(querySelectUserPackageByUserID, userid)
	if err != nil {
		fmt.Printf("Error Query: %v\n", err)
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	var ip uint32

	for rows.Next() {
		r := model.UserPackage{}
		err := rows.Scan(
			&r.ID, &r.UniqueID, &r.ParentID, &r.UserID, &r.PackageID, &r.StartDate, &r.StopDate, &r.Enabled, &r.ModifyDate,
			&ip, &r.IPNetwork, &r.MACAddress, &r.Description,
		)

		if err != nil {
			fmt.Printf("Error Scan: %v\n", err)
			return nil, store.ErrRecordNotFound
		}

		r.IPAddress = net.Ip2bytes(ip)

		got = append(got, r)
	}

	return got, nil
}
