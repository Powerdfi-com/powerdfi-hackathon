import { ReactNode } from "react"

const Modal = ({ children, onTapOutside }: { children: ReactNode, onTapOutside?: () => void }) => {
    return (
        <div className="z-50 bg-black-shade/40 backdrop-blur-md fixed top-0 left-0 bottom-0 right-0 flex items-center justify-center p-12" onClick={onTapOutside}>
            {children}
        </div>
    )
}

export default Modal