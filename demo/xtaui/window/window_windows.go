package window

//var (
//	dwmAPI                        = syscall.NewLazyDLL("dwmapi.dll")
//	_DwmExtendFrameIntoClientArea = dwmAPI.NewProc("DwmExtendFrameIntoClientArea")
//)
//
//type Margins struct {
//	CxLeftWidth, CxRightWidth, CyTopHeight, CyBottomHeight int32
//}
//
//func ExtendFrameIntoClientArea(hwnd uintptr, margins Margins) {
//	_, _, _ = _DwmExtendFrameIntoClientArea.Call(hwnd, uintptr(unsafe.Pointer(&margins)))
//}
//func (m *TMainWindow) wndProc(hwnd types.HWND, message uint32, wParam, lParam uintptr) uintptr {
//	switch message {
//	case messages.WM_DPICHANGED:
//		if !lcl.Application.Scaled() {
//			newWindowSize := (*types.TRect)(unsafe.Pointer(lParam))
//			win.SetWindowPos(m.Handle(), uintptr(0),
//				newWindowSize.Left, newWindowSize.Top, newWindowSize.Right-newWindowSize.Left, newWindowSize.Bottom-newWindowSize.Top,
//				win.SWP_NOZORDER|win.SWP_NOACTIVATE)
//		}
//	}
//	switch message {
//	case messages.WM_ACTIVATE:
//		ExtendFrameIntoClientArea(m.Handle(), Margins{CxLeftWidth: 1, CxRightWidth: 1, CyTopHeight: 1, CyBottomHeight: 1})
//	case messages.WM_NCCALCSIZE:
//		if wParam != 0 {
//			isMaximize := uint32(win.GetWindowLong(m.Handle(), win.GWL_STYLE))&win.WS_MAXIMIZE != 0
//			if isMaximize {
//				rect := (*types.TRect)(unsafe.Pointer(lParam))
//				monitor := winapi.MonitorFromRect(*rect, winapi.MONITOR_DEFAULTTONULL)
//				if monitor != 0 {
//					var monitorInfo types.TagMonitorInfo
//					monitorInfo.CbSize = types.DWORD(unsafe.Sizeof(monitorInfo))
//					if winapi.GetMonitorInfo(monitor, &monitorInfo) {
//						*rect = monitorInfo.RcWork
//					}
//				}
//			}
//			return 0
//		}
//	}
//
//	return win.CallWindowProc(m.oldWndPrc, uintptr(hwnd), message, wParam, lParam)
//}
//
//// 该函数调用可能会影响窗口的一些默认行为，需要知道在合适的时机调用它
//func (m *TMainWindow) _HookWndProcMessage() {
//	wndProcCallback := syscall.NewCallback(m.wndProc)
//	m.oldWndPrc = win.SetWindowLongPtr(m.Handle(), win.GWL_WNDPROC, wndProcCallback)
//}
//
//func (m *TMainWindow) _RestoreWndProc() {
//	if m.oldWndPrc != 0 {
//		win.SetWindowLongPtr(m.Handle(), win.GWL_WNDPROC, m.oldWndPrc)
//		m.oldWndPrc = 0
//	}
//}
