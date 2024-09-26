import { useState, useEffect } from "react";
import { format } from "date-fns";
import idLocale from "date-fns/locale/id";

export function RealtimeDate() {
  const [waktuSekarang, setWaktuSekarang] = useState(new Date());

  useEffect(() => {
    const intervalId = setInterval(() => {
      setWaktuSekarang(new Date());
    }, 1000);

    return () => clearInterval(intervalId);
  }, []);

  const tanggal = format(waktuSekarang, "EEEE, dd MMMM yyyy", {
    locale: idLocale,
  });

  return tanggal;
}

export function RealtimeClock() {
  const [waktuSekarang, setWaktuSekarang] = useState(new Date());

  useEffect(() => {
    const intervalId = setInterval(() => {
      setWaktuSekarang(new Date());
    }, 1000);

    return () => clearInterval(intervalId);
  }, []);

  const jam = format(waktuSekarang, "HH:mm:ss", {
    locale: idLocale,
  });

  return jam;
}
