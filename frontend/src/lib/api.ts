import axios from "axios";

const API = "http://localhost:8080/api";

export async function getTransactions() {
  const res = await axios.get(`${API}/transactions?status=PENDING`);
  return res.data;
}

export async function getRules() {
  const res = await axios.get(`${API}/rules`);
  return res.data;
}