import React from "react";

const DataItem = ({ count, label, bgColor, textColor }) => (
  <dl className={`bg-${bgColor}-50 rounded-lg flex flex-col items-center justify-center h-[78px]`}>
    <dt className={`w-8 h-8 rounded-full bg-${bgColor}-100 text-${textColor}-600 text-sm font-medium flex items-center justify-center mb-1`}>
      {count}
    </dt>
    <dd className={`text-${textColor}-600 text-sm font-medium`}>{label}</dd>
  </dl>
);

const Section = ({ span, cols, title, borderColor, children }) => (
  <div className={`${span} w-full h-fit bg-gray-50 border-b-4 border-${borderColor}-500 rounded-lg shadow p-4`}>
    <div className="flex justify-between mb-3">
      <div className="flex items-center">
        <div className="flex justify-center items-center">
          <h5 className="text-xl font-bold leading-none text-gray-900 pe-1">{title}</h5>
        </div>
      </div>
    </div>
    <div className="bg-white shadow p-3 rounded-lg">
      <div className={`grid ${cols} gap-3 mb-2 text-center`}>
        {children}
      </div>
    </div>
  </div>
);

export const FeatureDashboard = (props) => {
  const {
    sag,
    iso,
    memo,
    surat,
    beritaAcara,
    sk,
    project,
    rapat,
    perdin,
    cuti,
    masuk,
    keluar,
  } = props;

  return (
    <div className="grid grid-cols-4 gap-5">
      <Section span="col-span-3" cols="grid-cols-6" title="Data Dokumen" borderColor="blue">
        <DataItem count={sag.length} label="SAG" bgColor="teal" textColor="teal" />
        <DataItem count={iso.length} label="ISO" bgColor="cyan" textColor="cyan" />
        <DataItem count={memo.length} label="Memo" bgColor="sky" textColor="sky" />
        <DataItem count={surat.length} label="Surat" bgColor="blue" textColor="blue" />
        <DataItem count={beritaAcara.length} label="Berita Acara" bgColor="indigo" textColor="indigo" />
        <DataItem count={sk.length} label="SK" bgColor="indigo" textColor="violet" />
      </Section>

      <Section span="col-span-1" cols="grid-cols-1" title="Data Rencana Kerja" borderColor="green">
        <DataItem count={project.length} label="Project" bgColor="green" textColor="green" />
      </Section>

      <Section span="col-span-3" cols="grid-cols-3" title="Data Kegiatan & Proses" borderColor="amber">
        <DataItem count={rapat.length} label="Ruang Rapat" bgColor="orange" textColor="orange" />
        <DataItem count={perdin.length} label="Perjalanan Dinas" bgColor="amber" textColor="amber" />
        <DataItem count={cuti.length} label="Jadwal Cuti" bgColor="yellow" textColor="yellow" />
      </Section>

      <Section span="col-span-1" cols="grid-cols-2" title="Data & Informasi" borderColor="red">
        <DataItem count={masuk.length} label="Surat Masuk" bgColor="red" textColor="red" />
        <DataItem count={keluar.length} label="Surat Keluar" bgColor="orange" textColor="orange" />
      </Section>
    </div>
  );
};