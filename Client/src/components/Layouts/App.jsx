import { Avatar, Label, Sidebar } from "flowbite-react";
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
import { GetNotifRapat } from "../../../API/KegiatanProses/Notification/NotifRuangRapat";

const App = (props) => {
  const { services, children } = props;
  const [Notification, setNotification] = useState([]);
  const [isSidebarOpen, setIsSidebarOpen] = useState(
    services === "Ruang Rapat" || services === "Jadwal Cuti" ? false : true
  );
  const [userDetails, setUserDetails] = useState({ username: "", email: "" });
  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      const decoded = jwtDecode(token);
      setUserDetails({
        username: decoded.username,
        email: decoded.email,
        role: decoded.role,
      });
    }
  }, []);

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);

  const handleSignOut = async () => {
    try {
      // Panggil endpoint logout
      const response = await fetch("http://localhost:8080/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });

      // Cek status respons
      if (response.ok) {
        console.log("Logout berhasil");
        // Menghapus token dari localStorage
        localStorage.removeItem("token");
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
    GetNotifRapat((data) => {
      setNotification(data);
    });
  }, []);

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
                <Sidebar.Item icon={HiArrowSmRight} href="/sag">
                  SAG
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/iso">
                  ISO
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/memo">
                  MEMO
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/surat">
                  SURAT
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/berita-acara">
                  Berita Acara
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/sk">
                  SK
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
              <Sidebar.Collapse
                icon={HiOutlineMenuAlt2}
                label="Kegiatan & Proses"
              >
                <Sidebar.Item icon={HiArrowSmRight} href="/ruang-rapat">
                  Ruang Rapat
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/perjalanan-dinas">
                  Perjalanan Dinas
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/jadwal-cuti">
                  Jadwal Cuti
                </Sidebar.Item>
              </Sidebar.Collapse>
              <Sidebar.Collapse icon={HiOutlineServer} label="Data & Informasi">
                <Sidebar.Item icon={HiArrowSmRight} href="/surat-masuk">
                  Surat Masuk
                </Sidebar.Item>
                <Sidebar.Item icon={HiArrowSmRight} href="/surat-keluar">
                  Surat Keluar
                </Sidebar.Item>
              </Sidebar.Collapse>
              {userDetails.role === "admin" ? (
                <Sidebar.Item icon={AiOutlineUsergroupAdd} href="/register">
                  Tambah User
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
              }
            >
              <Dropdown.Header>
                <span className="block text-sm">Notification</span>
                <span className="block truncate text-sm font-medium">
                  {Notification.length} Jadwal!
                </span>
              </Dropdown.Header>
              {Notification.map((event) => {
                return (
                  <Dropdown.Item key={event.ID}>
                    <span className="block text-sm truncate">
                      Jadwal rapat {event.Message}
                    </span>
                  </Dropdown.Item>
                );
              })}
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
