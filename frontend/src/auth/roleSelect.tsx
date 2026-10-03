import {getKey} from "../util/language.ts";

export default function RoleSelect({ value, onChange }: {value: string, onChange: (role: string) => void}) {
    return (
        <select className="poppins" value={value} onChange={(e) => onChange(e.currentTarget.value)}>
            <option value="guardian">{getKey("ROLE.GUARDIAN")}</option>
            <option value="student">{getKey("ROLE.STUDENT")}</option>
            <option value="teacher">{getKey("ROLE.TEACHER")}</option>
        </select>
    )
}
