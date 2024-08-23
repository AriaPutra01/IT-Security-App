import { Avatar, Button, Label, Sidebar } from "flowbite-react";
import {
  HiChartPie,
  HiDocumentDuplicate,
  HiOutlineNewspaper,
  HiOutlineMenuAlt2,
  HiOutlineServer,
  HiArrowSmRight,
} from "react-icons/hi";
import { Dropdown } from "flowbite-react";
import { DarkThemeToggle } from "flowbite-react";
import { useState } from "react";

const App = (props) => {
  const { services, children } = props;
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);

  const handleSignOut = () => {
    // Menghapus token dari localStorage
    localStorage.removeItem('token');

    // Redirect ke halaman login
    window.location.href = '/login';
  };

  return (
    <div
      className={
        isSidebarOpen
          ? "grid h-screen grid-cols-2fr gap-6 p-4 dark:bg-gray-900"
          : "grid h-screen grid-cols-1 gap-6 p-4 dark:bg-gray-900"
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
              class="w-[30px] h-[30px] text-gray-800 dark:text-white cursor-pointer"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
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
            </Sidebar.ItemGroup>
          </Sidebar.Items>
        </Sidebar>
      ) : null}
      <div className="grid overflow-auto grid-rows-2fr w-full space-y-4">
        <div className="h-14 pb-2 flex justify-between border-b-2 border-gray-100 dark:border-gray-800">
          <div className="flex gap-2 items-end m-2">
            {isSidebarOpen ? null : (
              <svg
                onClick={toggleSidebar}
                class="w-[40px] h-[40px] text-gray-800 dark:text-white cursor-pointer"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-width="2"
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
            <DarkThemeToggle />
            <Dropdown
              arrowIcon={false}
              inline
              label={<Avatar status="online" rounded />}
            >
              <Dropdown.Header>
                <span className="block text-sm">Username</span>
                <span className="block truncate text-sm font-medium">
                  example@mail.com
                </span>
              </Dropdown.Header>
              <Dropdown.Item onClick={handleSignOut}>
                Sign out
              </Dropdown.Item>
            </Dropdown>
          </div>
        </div>
        {children}
      </div>
    </div>
  );
};
export default App;