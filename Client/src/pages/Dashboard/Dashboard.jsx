import App from "../../components/Layouts/App";

const token = localStorage.getItem('token'); // Ambil token dari localStorage
console.log("Token sent:", token);
const DashboardPage = () => {
  return <App services="dashboard"></App>;
};

export default DashboardPage;
