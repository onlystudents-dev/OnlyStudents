import {getKey} from "../util/language.ts";
import Select from "react-select";
import {selectStyles} from "../util/select.ts";

export default function RoleSelect({ value, onChange }: {value: string, onChange: (role: string) => void}) {
    const options = [
        { label: getKey("ROLE.GUARDIAN"), value: "guardian" },
        { label: getKey("ROLE.STUDENT"), value: "student" },
        { label: getKey("ROLE.TEACHER"), value: "teacher" },
    ]
    return <Select isSearchable={false} value={options.find(val => val.value === value)} options={options} onChange={val => val && onChange(val.value)} styles={selectStyles("100%", "3rem", "var(--hover-color)", "var(--card-color)", undefined, "0 0.5rem")} />
}
