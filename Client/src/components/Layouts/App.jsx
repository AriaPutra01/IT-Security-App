import { Avatar, Badge, Label } from "flowbite-react";
import { MdOutlineDashboard } from "react-icons/md";
import { HiOutlineClipboardDocumentList } from "react-icons/hi2";
import { GoProjectSymlink } from "react-icons/go";
import { GrPlan } from "react-icons/gr";
import { SlEnvolopeLetter } from "react-icons/sl";
import { FiUsers } from "react-icons/fi";
import { BiLogOut } from "react-icons/bi";
import { Dropdown } from "flowbite-react";
import { useState, useEffect } from "react";
import {
  deleteNotifRapat,
  GetNotifRapat,
} from "../../../API/KegiatanProses/Notification/NotifRuangRapat";
import { format } from "date-fns";
import { id } from "date-fns/locale";
import Swal from "sweetalert2";
import { useToken } from "../../context/TokenContext";
import Sidebar, { SidebarItem, SidebarCollapse } from "../Elements/Sidebar";

const App = (props) => {
  const { services, children } = props;
  const { token, userDetails } = useToken(); // Ambil token dari context
  const [Notification, setNotification] = useState([]);

  const handleSignOut = async () => {
    try {
      // Panggil endpoint logout
      const response = await fetch("http://localhost:8080/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        credentials: "include", // Sertakan cookie dalam permintaan
      });

      // Cek status respons
      if (response.ok) {
        // Redirect ke halaman login
        window.location.href = "/login";
      } else {
        const errorData = await response.json();
        alert("Logout gagal:", errorData);
      }
    } catch (error) {
      alert("Terjadi kesalahan saat melakukan logout:", error);
    }
  };

  // Fetch events
  useEffect(() => {
    GetNotifRapat((data) => {
      setNotification(data || []); // Pastikan setNotification menerima array
    });
  }, []);

  const handleDelete = async (id) => {
    Swal.fire({
      title: "Apakah Anda yakin?",
      text: "Anda akan menghapus Notif ini!",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    }).then(async (result) => {
      if (result.isConfirmed) {
        try {
          await deleteNotifRapat(id); // hapus data di API
          setNotification((prevData) =>
            prevData.filter((event) => event.id !== id) // Pastikan ini menggunakan 'id' yang benar
          );
        } catch (error) {
          Swal.fire("Gagal!", "Error saat hapus Notif:", error);
        }
      }
    });
  };

  return (
    <div className="grid grid-cols-2fr">
      <Sidebar
        img="../../../public/images/logobjb.png"
        title="Virtual Office"
        username={userDetails.username}
        email={userDetails.email}
      >
        <SidebarItem
          href="/dashboard"
          text="Dashboard"
          icon={<MdOutlineDashboard />}
        />
        <SidebarCollapse
          text="Dokumen"
          icon={<HiOutlineClipboardDocumentList />}
        >
          <SidebarItem href="/memo" text="Memo" />
          <SidebarItem href="/perjalanan-dinas" text="Perjalanan Dinas" />
        </SidebarCollapse>
        <SidebarCollapse text="Project" icon={<GoProjectSymlink />}>
          <SidebarItem href="/project" text="Project" />
          <SidebarItem href="/base-project" text="Base Project" />
        </SidebarCollapse>
        <SidebarCollapse text="Kegiatan" icon={<GrPlan />}>
          <SidebarItem href="/timeline-project" text="Timeline Project" />
          <SidebarItem href="/timeline-desktop" text="Timeline Desktop" />
          <SidebarItem href="/booking-rapat" text="Booking Ruang Rapat" />
          <SidebarItem href="/jadwal-rapat" text="Jadwal Rapat" />
          <SidebarItem href="/jadwal-cuti" text="Jadwal Cuti" />
        </SidebarCollapse>
        <SidebarCollapse text="Informasi" icon={<SlEnvolopeLetter />}>
          <SidebarItem href="/surat-masuk" text="Surat Masuk" />
          <SidebarItem href="/surat-keluar" text="Surat Keluar" />
        </SidebarCollapse>
        <SidebarItem href="/user" text="User" icon={<FiUsers />} />
        <SidebarItem
          onClick={handleSignOut}
          text="Logout"
          icon={<BiLogOut />}
        />
      </Sidebar>
      <div className="grid grid-rows-2fr">
        <div className="mx-4 pe-2 mt-2 flex justify-between border-b-2 border-gray-100">
          <div className="flex gap-2 items-end m-2">
            <div>
              <Label className="block text-sm">Halaman</Label>
              <Label className="block truncate text-sm font-medium ">
                <b className="uppercase">{services}</b>
              </Label>
            </div>
          </div>
          <Dropdown
            arrowIcon={false}
            inline
            label={
              <div className="relative">
                {Notification && Notification.length > 0 && (
                  <div className="absolute -translate-x-[3px] rounded-full bg-green-400">
                    <div className="w-full text-xs text-white px-[5px]">
                      {Notification.length}
                    </div>
                  </div>
                )}
                <svg
                  className="w-[34px] h-[34px] text-slate-700 dark:text-white"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M17.133 12.632v-1.8a5.406 5.406 0 0 1-4.154-5.262.955.955 0 0 0 .021-.106V3.1a1 1 0 0 0-2 0v2.364a.955.955 0 0 0 .021.106 5.406 5.406 0 0 0-4.154 5.262v1.8C6.867 15.018 5 15.614 5 16.807 5 17.4 5 18 5.538 18h12.924C19 18 19 17.4 19 16.807c0-1.193-1.867-1.789-1.867-4.175ZM8.823 19a3.453 3.453 0 0 0 6.354 0H8.823Z" />
                </svg>
              </div>
            }
          >
            <Dropdown.Header>
              <h1 className="text-base">Notification</h1>
            </Dropdown.Header>
            {Notification.length === 0 ? (
              <Dropdown.Item className="block text-sm text-gray-600">
                <Badge color="warning">
                  <h1>Tidak ada jadwal dalam wakut dekat</h1>
                </Badge>
              </Dropdown.Item>
            ) : (
              Notification.map((event) => {
                const formattedStart = format(event.start, "dd MMMM HH:mm", {
                  locale: id,
                });
                return (
                  <div
                    key={event.id}
                    className="flex items-center justify-between gap-2"
                  >
                    <Dropdown.Item className="flex gap-4">
                      <div className="grid grid-cols-2">
                        <span className="col-span-2 text-start ms-2 font-bold text-base truncate w-48">
                          {event.title}
                        </span>
                        <span className="col-span-2">
                          Pada Waktu {formattedStart}
                        </span>
                      </div>
                      <div
                        className="block text-sm truncate cursor-pointer hover:scale-110 text-red-600  rounded transition-all"
                        onClick={() => {
                          handleDelete(event.id);
                        }}
                      >
                        <svg
                          className="w-6 h-6"
                          aria-hidden="true"
                          xmlns="http://www.w3.org/2000/svg"
                          width="24"
                          height="24"
                          fill="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            fillRule="evenodd"
                            d="M8.586 2.586A2 2 0 0 1 10 2h4a2 2 0 0 1 2 2v2h3a1 1 0 1 1 0 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V8a1 1 0 0 1 0-2h3V4a2 2 0 0 1 .586-1.414ZM10 6h4V4h-4v2Zm1 4a1 1 0 1 0-2 0v8a1 1 0 1 0 2 0v-8Zm4 0a1 1 0 1 0-2 0v8a1 1 0 1 0 2 0v-8Z"
                            clipRule="evenodd"
                          />
                        </svg>
                      </div>
                    </Dropdown.Item>
                  </div>
                );
              })
            )}
          </Dropdown>
        </div>
        <div className="mt-4 px-2 w-full overflow-auto">{children}</div>
      </div>
    </div>
  );
};
export default App;
