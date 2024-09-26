import { useToken } from "../../context/TokenContext";
import { Link, Navigate } from "react-router-dom";
import "../../welcome.css";
export function WelcomePage() {
  const { token } = useToken();
  if (token) return <Navigate to="/dashboard" />;
  return (
    <section className="h-screen flex justify-center items-center">
      <div className="text-center">
        <img
          className="mx-auto mb-8 w-28"
          src="/public/images/logobjbputih.png"
          alt="bjb"
        />
        <h1 className="mb-4 text-4xl font-extrabold tracking-tight leading-none text-slate-100 md:text-4xl lg:text-5xl ">
          Selamat Datang di Dashboard ITS
        </h1>
        <p className="mb-8 text-lg font-normal text-slate-100 lg:text-xl sm:px-16 xl:px-48 ">
          Aplikasi ini digunakan untuk Mempermudah dalam Mengelola Proyek dan
          Control Data untuk Keamanan IT Security
        </p>
        <div className="flex flex-col items-center mb-8 lg:mb-16 space-y-4 sm:flex-row sm:justify-center sm:space-y-0 sm:space-x-4">
          <Link
            to="/login"
            className="inline-flex justify-center items-center py-3 px-5 text-base font-medium text-center bg-transparent border-2 hover:text-sky-400 hover:border-sky-400 hover:scale-105 transition-all text-white rounded-lg bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:ring-primary-300"
          >
            Mulai Gunakan
            <svg
              className="ml-2 -mr-1 w-5 h-5"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z"></path>
            </svg>
          </Link>
        </div>
      </div>
    </section>
  );
}
