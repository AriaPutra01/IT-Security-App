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
import { DarkThemeToggle } from "flowbite-react";

const App = (props) => {
  const { services, children } = props;
  return (
    <div className="grid h-screen grid-cols-2fr gap-6 p-4 dark:bg-gray-900">
      <Sidebar className="rounded-xl overflow-auto">
        <Sidebar.Logo
          href="/dashboard"
          img="/public/images/logobjb.png"
          imgAlt="Bank bjb logo"
        >
          Bank bjb
        </Sidebar.Logo>
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
      <div className="grid overflow-auto grid-rows-2fr w-full space-y-4">
        <div className="h-14 px-4 flex justify-between border-b-2 border-gray-100 dark:border-gray-800">
          <div>
            <Label className="block text-sm">Halaman</Label>
            <Label className="block truncate text-sm font-medium ">
              <b className="uppercase">{services}</b>
            </Label>
          </div>
          <div className="flex items-center gap-4">
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
              <Dropdown.Item>
                <a className="block w-full" href="#">
                  Sign out
                </a>
              </Dropdown.Item>
            </Dropdown>
          </div>
        </div>
        <div>{children}</div>
      </div>
    </div>
  );
};
export default App;
