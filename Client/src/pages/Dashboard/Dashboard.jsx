import App from "../../components/Layouts/App";

const token = localStorage.getItem('token'); // Ambil token dari localStorage
const DashboardPage = () => {
  return <App services="dashboard"></App>;
};

export default DashboardPage;
