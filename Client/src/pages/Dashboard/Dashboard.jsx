import { useEffect, useState } from "react";
import { FeatureDashboard } from "../../components/Fragments/Dashboard/FeatureDashboard";
import App from "../../components/Layouts/App";
import { getSuratKeluars } from "../../../API/DataInformasi/SuratKeluar.service";
import { getSuratMasuks } from "../../../API/DataInformasi/SuratMasuk.service";
import { getMemos } from "../../../API/Dokumen/memo.service";
import { getPerdins } from "../../../API/Dokumen/PerjalananDinas.service";
import { getEventsProject } from "../../../API/KegiatanProses/TimelineProject.service";
import { getEventsDesktop } from "../../../API/KegiatanProses/TimelineDesktop.service";
import { getBookingRapat } from "../../../API/KegiatanProses/BookingRapat.service";
import { getCutis } from "../../../API/KegiatanProses/JadwalCuti.service";
import { getRapats } from "../../../API/KegiatanProses/JadwalRapat.service";
import { getProjects } from "../../../API/RencanaKerja/Project.service";
import { useToken } from "../../context/TokenContext";

const useFetchData = (fetchFunction) => {
  const [data, setData] = useState([]);
  useEffect(() => {
    fetchFunction(setData);
  }, [fetchFunction]);
  return data;
};

export const DashboardPage = () => {
  const { token } = useToken(); // Ambil token dari context

  const SuratKeluar = useFetchData(getSuratKeluars);
  const SuratMasuk = useFetchData(getSuratMasuks);
  const Memo = useFetchData(getMemos);
  const Perdin = useFetchData(getPerdins);
  const TimelineProject = useFetchData(getEventsProject);
  const TimelineDesktop = useFetchData(getEventsDesktop);
  const Booking = useFetchData(getBookingRapat);
  const Cuti = useFetchData(getCutis);
  const Rapat = useFetchData(getRapats);
  const Project = useFetchData(getProjects);

  return (
    <App services="Dashboard">
      <FeatureDashboard
        memo={Memo}
        perdin={Perdin}
        project={Project}
        timelineProject={TimelineProject}
        timelineDesktop={TimelineDesktop}
        booking={Booking}
        rapat={Rapat}
        cuti={Cuti}
        masuk={SuratMasuk}
        keluar={SuratKeluar}
      />
    </App>
  );
};
