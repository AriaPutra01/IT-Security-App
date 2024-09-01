import { useEffect, useState } from "react";
import { FeatureDashboard } from "../../components/Fragments/Dashboard/FeatureDashboard";
import App from "../../components/Layouts/App";
import { getSuratKeluars } from "../../../API/DataInformasi/SuratKeluar.service";
import { getSuratMasuks } from "../../../API/DataInformasi/SuratMasuk.service";
import { getBeritaAcaras } from "../../../API/Dokumen/BeritaAcara.service";
import { getIsos } from "../../../API/Dokumen/iso.service";
import { getMemos } from "../../../API/Dokumen/memo.service";
import { getSags } from "../../../API/Dokumen/sag.service";
import { getSks } from "../../../API/Dokumen/Sk.service";
import { getSurats } from "../../../API/Dokumen/surat.service";
import { getCutis } from "../../../API/KegiatanProses/JadwalCuti.service";
import { getPerdins } from "../../../API/KegiatanProses/PerjalananDinas.service";
import { getRapats } from "../../../API/KegiatanProses/RuangRapat.service";
import { getProjects } from "../../../API/RencanaKerja/Project.service";

const useFetchData = (fetchFunction) => {
  const [data, setData] = useState([]);
  useEffect(() => {
    fetchFunction(setData);
  }, [fetchFunction]);
  return data;
};

export const DashboardPage = () => {
  const token = localStorage.getItem("token"); // Ambil token dari localStorage

  const SuratKeluar = useFetchData(getSuratKeluars);
  const SuratMasuk = useFetchData(getSuratMasuks);
  const BeritaAcara = useFetchData(getBeritaAcaras);
  const Iso = useFetchData(getIsos);
  const Memo = useFetchData(getMemos);
  const Sag = useFetchData(getSags);
  const Sk = useFetchData(getSks);
  const Surat = useFetchData(getSurats);
  const Cuti = useFetchData(getCutis);
  const Perdin = useFetchData(getPerdins);
  const Rapat = useFetchData(getRapats);
  const Project = useFetchData(getProjects);

  return (
    <App services="Dashboard">
      <FeatureDashboard
        sag={Sag}
        iso={Iso}
        memo={Memo}
        surat={Surat}
        beritaAcara={BeritaAcara}
        sk={Sk}
        project={Project}
        rapat={Rapat}
        perdin={Perdin}
        cuti={Cuti}
        masuk={SuratMasuk}
        keluar={SuratKeluar}
      />
    </App>
  );
};
