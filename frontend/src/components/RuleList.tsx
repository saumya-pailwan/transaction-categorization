interface Rule {
  id: string;
  pattern_type: string;
  pattern_value: string;
  category_id: string;
}

interface RuleListProps {
  rules: Rule[];
}

export default function RuleList({ rules }: RuleListProps) {
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