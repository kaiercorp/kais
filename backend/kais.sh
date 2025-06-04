#!/bin/bash
#
################################################################
APP=main
# main 실행파일과 같은 디렉토리에 PID 파일 생성
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MAIN_PATH="$SCRIPT_DIR/$APP"
PID="$SCRIPT_DIR/$APP.pid"
# LOG는 start 함수에서 동적으로 설정됨
LOG_DIR="$SCRIPT_DIR/log"
COMMAND="$MAIN_PATH -p $PID"
################################################################

start() {
echo "==== Start KAI.S ===="

# LOG 파일을 시작 시점의 날짜로 설정
LOG="$LOG_DIR/$(date '+%Y%m%d').log"
echo "Log file: $LOG"

# 1. 실행 파일 존재 확인
if [ ! -f "$MAIN_PATH" ]; then
    echo "Error: Executable file not found: $MAIN_PATH"
    return 1
fi

if [ ! -x "$MAIN_PATH" ]; then
    echo "Error: File is not executable: $MAIN_PATH"
    return 1
fi

# 2. PID 파일 확인 및 기존 프로세스 체크
if [ -f "$PID" ]; then
    PID_NUM=$(cat "$PID" 2>/dev/null)
    if [ -n "$PID_NUM" ] && ps -p "$PID_NUM" > /dev/null 2>&1; then
        echo "Already running with PID: $PID_NUM"
        echo "Use 'stop' first, or 'restart' to restart the service."
        return 1
    else
        echo "Found stale PID file. Removing..."
        rm -f "$PID"
    fi
fi

# 3. 이미 실행 중인 프로세스가 있는지 확인 (PID 파일 없이 실행 중인 경우)
EXISTING_PIDS=$(pgrep -f "$MAIN_PATH" 2>/dev/null)
if [ -n "$EXISTING_PIDS" ]; then
    echo "Warning: Found running processes without PID file:"
    echo "$EXISTING_PIDS"
    echo "Please stop existing processes first."
    return 1
fi

# 4. 로그 디렉토리 생성
if [ ! -d "$LOG_DIR" ]; then
    echo "Creating log directory: $LOG_DIR"
    mkdir -p "$LOG_DIR" || {
        echo "Error: Cannot create log directory: $LOG_DIR"
        return 1
    }
fi

# 5. PID 파일 디렉토리 확인
PID_DIR=$(dirname "$PID")
if [ ! -d "$PID_DIR" ]; then
    echo "Creating PID directory: $PID_DIR"
    mkdir -p "$PID_DIR" || {
        echo "Error: Cannot create PID directory: $PID_DIR"
        return 1
    }
fi

# 6. 프로세스 시작
echo "Starting $APP..."
echo "Command: $COMMAND"

# nohup로 백그라운드 실행
if nohup $COMMAND > /dev/null 2>&1 & then
    PROCESS_PID=$!
    echo "$PROCESS_PID" > "$PID"
    
    # 짧은 시간 대기 후 프로세스가 실제로 시작되었는지 확인
    sleep 2
    
    if ps -p "$PROCESS_PID" > /dev/null 2>&1; then
        echo "Successfully started with PID: $PROCESS_PID"
        echo "PID file: $PID"
        echo "$(date '+%Y-%m-%d %X'): START (PID: $PROCESS_PID)" >> "$LOG"
        return 0
    else
        echo "Error: Process failed to start or exited immediately"
        echo "Check the application logs for details"
        rm -f "$PID"
        return 1
    fi
else
    echo "Error: Failed to execute command"
    rm -f "$PID"
    return 1
fi
}
stop() {
echo "==== Stop KAI.S ===="

# stop 시에도 현재 날짜의 로그 파일 사용
LOG="$LOG_DIR/$(date '+%Y%m%d').log"

if [ -f "$PID" ]; then
    PID_NUM=$(cat "$PID")
    echo "Found PID file with PID: $PID_NUM"
    
    # PID가 실제로 실행 중인지 확인
    if ps -p "$PID_NUM" > /dev/null 2>&1; then
        echo "Process is running. Attempting graceful shutdown..."
        
        # 1단계: SIGTERM으로 정상 종료 시도 (15초 대기)
        if kill -TERM "$PID_NUM" 2>/dev/null; then
            echo "Sent SIGTERM signal. Waiting for graceful shutdown..."
            for i in {1..15}; do
                if ! ps -p "$PID_NUM" > /dev/null 2>&1; then
                    echo "Process terminated gracefully."
                    break
                fi
                sleep 1
                echo -n "."
            done
            echo
        fi
        
        # 2단계: 여전히 실행 중이면 SIGKILL로 강제 종료
        if ps -p "$PID_NUM" > /dev/null 2>&1; then
            echo "Process still running. Force killing..."
            if kill -KILL "$PID_NUM" 2>/dev/null; then
                sleep 2
                if ps -p "$PID_NUM" > /dev/null 2>&1; then
                    echo "Warning: Process may still be running"
                else
                    echo "Process force killed."
                fi
            else
                echo "Failed to kill process"
            fi
        fi
    else
        echo "Process not running (PID $PID_NUM not found)"
    fi
    
    # PID 파일 제거
    rm -f "$PID"
    echo "PID file removed."
    echo "$(date '+%Y-%m-%d %X'): STOP" >> "$LOG"
    
else
    echo "No PID file found. Checking for running processes..."
    
    # PID 파일이 없는 경우 프로세스명으로 검색하여 정리
    PIDS=$(pgrep -f "$MAIN_PATH" 2>/dev/null)
    if [ -n "$PIDS" ]; then
        echo "Found running processes: $PIDS"
        echo "Cleaning up orphaned processes..."
        
        # 정상 종료 시도
        echo "$PIDS" | xargs kill -TERM 2>/dev/null
        sleep 3
        
        # 강제 종료
        REMAINING=$(pgrep -f "$MAIN_PATH" 2>/dev/null)
        if [ -n "$REMAINING" ]; then
            echo "Force killing remaining processes: $REMAINING"
            echo "$REMAINING" | xargs kill -KILL 2>/dev/null
        fi
        
        echo "Cleanup completed."
        echo "$(date '+%Y-%m-%d %X'): STOP (orphaned processes)" >> "$LOG"
    else
        echo "No running processes found."
    fi
fi

echo "Stop completed."
}

case "$1" in
    'start')
            start
            ;;
    'stop')
            stop
            ;;
    'restart')
            stop ; echo "Sleeping..."; sleep 1 ;
            start
            ;;
    *)
            echo
            echo "Usage: $0 { start | stop | restart }"
            echo
            exit 1
            ;;
esac

exit 0