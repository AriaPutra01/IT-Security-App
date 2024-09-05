import { Avatar, Badge, Label, Sidebar } from "flowbite-react";
import {
  HiChartPie,
  HiDocumentDuplicate,
  HiOutlineNewspaper,
  HiOutlineMenuAlt2,
  HiOutlineServer,
  HiArrowSmRight,
} from "react-icons/hi";
import { Dropdown } from "flowbite-react";
import { useState, useEffect } from "react";
import { jwtDecode } from "jwt-decode";
import { AiOutlineUsergroupAdd } from "react-icons/ai";
import {
  deleteNotifRapat,
  GetNotifRapat,
} from "../../../API/KegiatanProses/Notification/NotifRuangRapat";
import { format } from "date-fns";
import { id } from "date-fns/locale";
import Swal from "sweetalert2";
import { useToken } from "../../context/TokenContext";

const App = (props) => {
  const { services, children } = props;
  const { token, userDetails } = useToken(); // Ambil token dari context
  const [Notification, setNotification] = useState([]);
  const [isSidebarOpen, setIsSidebarOpen] = useState(
    services === "Timeline Desktop" ||
      services === "Booking Rapat" ||
      services === "Jadwal Rapat" ||
      services === "Jadwal Cuti"
      ? false
      : true
  );

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);

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
        console.error("Logout gagal:", errorData);
      }
    } catch (error) {
      console.error("Terjadi kesalahan saat melakukan logout:", error);
    }
  };

  // Fetch events
  useEffect(() => {
    GetNotifRapat((event) => {
      setNotification(event);
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
            prevData.filter((event) => event.ID !== id)
          );
        } catch (error) {
          Swal.fire("Gagal!", "Error saat hapus Notif:", error);
        }
      }
    });
  };

  return (
    <div
      className={
        isSidebarOpen
          ? "grid h-screen grid-cols-2fr gap-6 p-4"
          : "grid h-screen grid-cols-1 gap-6 p-4"
      }
    >
      {isSidebarOpen ? (
        <Sidebar className="rounded-xl h-full overflow-auto">
          <div className="flex justify-between mr-2">
            <Sidebar.Logo
              href="/dashboard"
              img="/public/images/logobjb.png"
              imgAlt="Bank bjb logo"
            >
              Bank bjb
            </Sidebar.Logo>
            <svg
              onClick={toggleSidebar}
              className="w-[30px] h-[30px] text-gray-800 cursor-pointer"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M6 18 17.94 6M18 18 6.06 6"
              />
            </svg>
          </div>
          <Sidebar.Items>
            <Sidebar.ItemGroup>
              <Sidebar.Item href="/dashboard" icon={HiChartPie}>
                Dashboard
              </Sidebar.Item>
              <Sidebar.Collapse icon={HiDocumentDuplicate} label="Dokumen">
                <Sidebar.Item icon={HiArrowSmRight} href="/memo">
                  MEMO
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/perjalanan-dinas">
                  Perjalanan Dinas
                </Sidebar.Item>
              </Sidebar.Collapse>
              <Sidebar.Collapse icon={HiOutlineNewspaper} label="Rencana Kerja">
                <Sidebar.Item icon={HiArrowSmRight} href="/project">
                  Project
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/base-project">
                  Base Project
                </Sidebar.Item>
              </Sidebar.Collapse>
              <Sidebar.Collapse icon={HiOutlineMenuAlt2} label="Kegiatan">
                <Sidebar.Item icon={HiArrowSmRight} href="/timeline">
                  Timeline Desktop
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/booking-rapat">
                  Booking Rapat
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/jadwal-rapat">
                  Jadwal Rapat
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/jadwal-cuti">
                  Jadwal Cuti
                </Sidebar.Item>
              </Sidebar.Collapse>
              <Sidebar.Collapse icon={HiOutlineServer} label="Informasi">
                <Sidebar.Item icon={HiArrowSmRight} href="/surat-masuk">
                  Surat Masuk
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/surat-keluar">
                  Surat Keluar
                </Sidebar.Item>
              </Sidebar.Collapse>
              {userDetails.role === "admin" ? (
                <Sidebar.Item icon={AiOutlineUsergroupAdd} href="/user">
                  User Management
                </Sidebar.Item>
              ) : null}
            </Sidebar.ItemGroup>
          </Sidebar.Items>
        </Sidebar>
      ) : null}
      <div className="grid grid-rows-2fr w-full space-y-4">
        <div className="h-14 pb-2 flex justify-between border-b-2 border-gray-100">
          <div className="flex gap-2 items-end m-2">
            {isSidebarOpen ? null : (
              <svg
                onClick={toggleSidebar}
                className="w-[40px] h-[40px] text-gray-800 cursor-pointer"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeWidth="2"
                  d="M5 7h14M5 12h14M5 17h14"
                />
              </svg>
            )}
            <div>
              <Label className="block text-sm">Halaman</Label>
              <Label className="block truncate text-sm font-medium ">
                <b className="uppercase">{services}</b>
              </Label>
            </div>
          </div>
          <div className="flex items-center gap-4 ">
            <Dropdown
              arrowIcon={false}
              inline
              label={
                <div className="relative">
                  {Notification.length > 0 && (
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
                    <path d="M17.133 12.632v-1.8a5.406 5.406 0 0 0-4.154-5.262.955.955 0 0 0 .021-.106V3.1a1 1 0 0 0-2 0v2.364a.955.955 0 0 0 .021.106 5.406 5.406 0 0 0-4.154 5.262v1.8C6.867 15.018 5 15.614 5 16.807 5 17.4 5 18 5.538 18h12.924C19 18 19 17.4 19 16.807c0-1.193-1.867-1.789-1.867-4.175ZM8.823 19a3.453 3.453 0 0 0 6.354 0H8.823Z" />
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
                      key={event.ID}
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
                            handleDelete(event.ID);
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
            <Dropdown
              arrowIcon={false}
              inline
              label={<Avatar status="online" rounded />}
            >
              <Dropdown.Header>
                <span className="block text-sm">{userDetails.username}</span>
                <span className="block truncate text-sm font-medium">
                  {userDetails.email}
                </span>
              </Dropdown.Header>
              <Dropdown.Item onClick={handleSignOut}>Sign out</Dropdown.Item>
            </Dropdown>
          </div>
        </div>
        {children}
      </div>
    </div>
  );
};
export default App;
