export default function SubjectIcon({ icon }: { icon: string }) {
    return (
        <>
            {icon && <div className="material text-3xl flex items-center pr-2">{icon}</div>}
        </>
    )
}