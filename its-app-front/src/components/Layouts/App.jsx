"use client";
import { Avatar, Label, Sidebar, TextInput } from "flowbite-react";
import {
  HiChartPie,
  HiDocumentDuplicate,
  HiOutlineNewspaper,
  HiOutlineMenuAlt2,
  HiOutlineServer,
  HiArrowSmRight,
} from "react-icons/hi";
import { Dropdown } from "flowbite-react";
import { DarkThemeToggle, Flowbite } from "flowbite-react";
// Pages
// Dashboard
import FeatureDashboard from "../Fragments/Dashboard/FeatureDashboard";
// Dokumen
import { IndexSag } from "../Fragments/Services/Dokumen/CrudSag";
// Rencana Kerja
import { IndexProject } from "../Fragments/Services/RencanaKerja/CrudProject";
import { IndexBaseProject } from "../Fragments/Services/RencanaKerja/CrudBaseProject";
// Kegiatan Proses
import { IndexJadwalCuti } from "../Fragments/Services/KegiatanProses/CrudJadwalCuti";
import { IndexPerjalananDinas } from "../Fragments/Services/KegiatanProses/CrudPerjalananDinas";
import { IndexRuangRapat } from "../Fragments/Services/KegiatanProses/CrudRuangRapat";
// Data Informasi
import { IndexSuratMasuk } from "../Fragments/Services/DataInformasi/CrudSuratMasuk";
import { IndexSuratKeluar } from "../Fragments/Services/DataInformasi/CrudSuratKeluar";

const App = (props) => {
  const { services } = props;
  return (
    <Flowbite>
      <div className="flex dark:bg-gray-900">
        <Sidebar className="m-4 rounded-full h-screen">
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
        <div className="flex flex-col mt-4 w-full space-y-4">
          <div className="h-14 mx-6 pb-4 flex justify-between border-b-2 border-gray-100 dark:border-gray-800">
            <div>
              <Label className="block text-sm">Halaman</Label>
              <Label className="block truncate text-sm font-medium ">
                <b className="uppercase">{services}</b>
              </Label>
            </div>
            <TextInput
              className="w-1/3"
              id="search"
              type="search"
              placeholder="Search"
            ></TextInput>
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
                  <a className="block w-full h-full" href="#">
                    Sign out
                  </a>
                </Dropdown.Item>
              </Dropdown>
            </div>
          </div>
          <div className="h-full m-4">
            {/* Dashboard */}
            {services === "dashboard" && <FeatureDashboard />}
            {/* Dokumen */}
            {services === "sag" && <IndexSag />}
            {/* Rencana Kerja */}
            {services === "project" && <IndexProject />}
            {services === "base-project" && <IndexBaseProject />}
            {/* Kegiatan dan Proses */}
            {services === "ruang-rapat" && <IndexRuangRapat />}
            {services === "perjalanan-dinas" && <IndexPerjalananDinas />}
            {services === "jadwal-cuti" && <IndexJadwalCuti />}
            {/* Data dan Inasi */}
            {services === "surat-masuk" && <IndexSuratMasuk />}
            {services === "surat-keluar" && <IndexSuratKeluar />}
          </div>
        </div>
      </div>
    </Flowbite>
  );
};
export default App;
