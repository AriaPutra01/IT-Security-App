import { useEffect, useState } from "react";
import { FeatureDashboard } from "../../components/Fragments/Dashboard/FeatureDashboard";
import App from "../../components/Layouts/App";
import { getSuratKeluars } from "../../../API/DataInformasi/SuratKeluar.service";
import { getSuratMasuks } from "../../../API/DataInformasi/SuratMasuk.service";
import { getMemos } from "../../../API/Dokumen/MemoSag.service";
import { getBeritaAcaras } from "../../../API/Dokumen/BeritaAcara.service";
import { getSurats } from "../../../API/Dokumen/Surat.service";
import { getSks } from "../../../API/Dokumen/Sk.service";
import { getPerdins } from "../../../API/Dokumen/PerjalananDinas.service";
import { getEventsProject } from "../../../API/KegiatanProses/TimelineProject.service";
import { getEventsDesktop } from "../../../API/KegiatanProses/TimelineDesktop.service";
import { getBookingRapat } from "../../../API/KegiatanProses/BookingRapat.service";
import { getMeetings } from "../../../API/KegiatanProses/Meeting.service";
import { getMeetingList } from "../../../API/KegiatanProses/MeetingSchedule.service";
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
  const Memos = useFetchData(getMemos);
  const BeritaAcara = useFetchData(getBeritaAcaras);
  const Surat = useFetchData(getSurats);
  const Sk = useFetchData(getSks);
  const Perdin = useFetchData(getPerdins);
  const TimelineProject = useFetchData(getEventsProject);
  const TimelineDesktop = useFetchData(getEventsDesktop);
  const Booking = useFetchData(getBookingRapat);
  const Meeting = useFetchData(getMeetings);
  const MeetingList = useFetchData(getMeetingList);
  const Cuti = useFetchData(getCutis);
  const Rapat = useFetchData(getRapats);
  const Project = useFetchData(getProjects);

  return (
    <App services="Dashboard">
      <FeatureDashboard
        memo={Memos}
        beritaAcara={BeritaAcara}
        surat={Surat}
        sk={Sk}
        perdin={Perdin}
        project={Project}
        timelineProject={TimelineProject}
        timelineDesktop={TimelineDesktop}
        meeting={Meeting}
        meetingList={MeetingList}
        booking={Booking}
        rapat={Rapat}
        cuti={Cuti}
        masuk={SuratMasuk}
        keluar={SuratKeluar}
      />
    </App>
  );
};
