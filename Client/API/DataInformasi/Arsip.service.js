import axios from 'axios';

const API_URL = 'http://localhost:8080/Arsip'; // Ganti dengan URL API Anda

export const getArsip = async (callback) => {
    try {
        const response = await axios.get(`${API_URL}`); // Update endpoint
        callback(response.data);
    } catch (error) {
        console.error("Error fetching arsip:", error);
    }
};

export const addArsip = async (data) => {
    try {
        await axios.post(`${API_URL}`, data); // Update endpoint
    } catch (error) {
        console.error("Error adding arsip:", error);
        throw error; // Melempar error untuk ditangani di komponen
    }
};

export const deleteArsip = async (id) => {
    try {
        await axios.delete(`${API_URL}/arsip/${id}`); // Update endpoint
    } catch (error) {
        console.error("Error deleting arsip:", error);
        throw error; // Melempar error untuk ditangani di komponen
    }
};

export const updateArsip = async (id, data) => {
    try {
        await axios.put(`${API_URL}/arsip/${id}`, data); // Update endpoint
    } catch (error) {
        console.error("Error updating arsip:", error);
        throw error; // Melempar error untuk ditangani di komponen
    }
};




