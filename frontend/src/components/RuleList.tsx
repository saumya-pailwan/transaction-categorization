export default function RuleList({ rules }) {
  return (
    <ul className="list-disc ml-6">
      {rules.map((rule) => (
        <li key={rule.id} className="mb-2">
          <strong>{rule.pattern_type}</strong>: "{rule.pattern_value}" → {rule.category_id}
        </li>
      ))}
    </ul>
  );
}