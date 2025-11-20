"use client";
import { useEffect, useState } from "react";
import { getRules } from "@/lib/api";
import RuleList from "@/components/RuleList";

export default function RulesPage() {
  const [rules, setRules] = useState([]);

  useEffect(() => {
    async function load() {
      const data = await getRules();
      setRules(data || []);
    }
    load();
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-semibold mb-4">Rules</h2>
      <RuleList rules={rules} />
    </div>
  );
}